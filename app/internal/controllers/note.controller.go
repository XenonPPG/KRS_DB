package controllers

import (
	"DB/internal/initializers"
	"context"

	desc "github.com/XenonPPG/KRS_CONTRACTS/gen/note_v1"
	"github.com/XenonPPG/KRS_CONTRACTS/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) CreateNote(ctx context.Context, req *desc.CreateNoteRequest) (*desc.Note, error) {
	note := &models.Note{
		Title:   req.GetTitle(),
		Content: req.GetContent(),
		UserID:  req.GetUserID(),
	}

	result := initializers.DB.Create(note)

	if result.Error != nil {
		return nil, result.Error
	}

	return &desc.Note{
		Id:        note.ID,
		Title:     note.Title,
		Content:   note.Content,
		UserID:    note.UserID,
		CreatedAt: timestamppb.New(note.CreatedAt),
	}, nil
}

func (s *Server) GetAllNotes(ctx context.Context, req *desc.GetAllNotesRequest) (*desc.GetAllNotesResponse, error) {
	var notes []models.Note

	// results with pagination
	result := initializers.DB.
		Where("user_id = ?", req.GetUserID()).
		Limit(int(req.Limit)).
		Offset(int(req.Offset)).
		Find(&notes)

	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch notes: %v", result.Error)
	}

	// total count of user's notes
	var notesAmount int64
	result = initializers.DB.Model(&models.Note{}).Where("user_id = ?", req.GetUserID()).Count(&notesAmount)
	if result.Error != nil {
		return nil, result.Error
	}

	// format models.Note to proto.Note
	protoNotes := make([]*desc.Note, 0, notesAmount)
	for _, n := range notes {
		protoNotes = append(protoNotes, &desc.Note{
			Id:        n.ID,
			Title:     n.Title,
			Content:   n.Content,
			UserID:    n.UserID,
			CreatedAt: timestamppb.New(n.CreatedAt),
		})
	}

	return &desc.GetAllNotesResponse{
		Notes:      protoNotes,
		TotalCount: int32(notesAmount),
	}, nil
}

func (s *Server) GetNote(ctx context.Context, req *desc.GetNoteRequest) (*desc.Note, error) {
	var note models.Note

	result := initializers.DB.First(&note, req.GetId())
	if result.Error != nil {
		return nil, result.Error
	}

	// check if note is owned by user
	if note.UserID != req.GetUserID() {
		return nil, status.Errorf(codes.PermissionDenied, "note is not owned by user")
	}

	return &desc.Note{
		Id:        note.ID,
		Title:     note.Title,
		Content:   note.Content,
		UserID:    note.UserID,
		CreatedAt: timestamppb.New(note.CreatedAt),
	}, nil
}

func (s *Server) UpdateNote(ctx context.Context, req *desc.UpdateNoteRequest) (*desc.Note, error) {
	oldNote, err := s.GetNote(ctx, &desc.GetNoteRequest{
		Id:     req.GetId(),
		UserID: req.GetUserID(),
	})
	if err != nil {
		return nil, err
	}
	if oldNote == nil {
		return nil, status.Errorf(codes.NotFound, "note not found")
	}

	// check if note is owned by user
	if oldNote.UserID != req.GetUserID() {
		return nil, status.Errorf(codes.PermissionDenied, "note is not owned by user")
	}

	// check if something was updated
	if oldNote.Title == req.GetTitle() && oldNote.Content == req.GetContent() {
		return nil, status.Errorf(codes.InvalidArgument, "nothing to update")
	}

	// update
	newNote := &models.Note{
		ID:        oldNote.Id,
		Title:     req.GetTitle(),
		Content:   req.GetContent(),
		UserID:    oldNote.UserID,
		CreatedAt: oldNote.CreatedAt.AsTime(),
	}
	result := initializers.DB.Model(&models.Note{ID: req.GetId()}).Updates(newNote)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "failed to update note: %v", result.Error)
	}

	return &desc.Note{
		Id:        newNote.ID,
		Title:     newNote.Title,
		Content:   newNote.Content,
		UserID:    newNote.UserID,
		CreatedAt: timestamppb.New(newNote.CreatedAt),
	}, nil
}

func (s *Server) DeleteNote(ctx context.Context, req *desc.DeleteNoteRequest) (*emptypb.Empty, error) {
	// check if note is owned by user
	note, err := s.GetNote(ctx, &desc.GetNoteRequest{
		Id:     req.GetId(),
		UserID: req.GetUserID(),
	})
	if err != nil {
		return nil, err
	}
	if note == nil {
		return nil, status.Errorf(codes.NotFound, "note not found")
	}

	if note.UserID != req.GetUserID() {
		return nil, status.Errorf(codes.PermissionDenied, "note is not owned by user")
	}

	result := initializers.DB.Delete(&models.Note{}, req.GetId())

	if result.Error != nil {
		return nil, result.Error
	}

	return &emptypb.Empty{}, nil
}
