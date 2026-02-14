package controllers

import (
	"DB/internal/initializers"
	"context"

	desc "github.com/XenonPPG/KRS_CONTRACTS/gen/note_v1"
	"github.com/XenonPPG/KRS_CONTRACTS/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
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
		Id:      note.ID,
		Title:   note.Title,
		Content: note.Content,
	}, nil
}

func (s *Server) GetAllNotes(ctx context.Context, req *desc.GetAllNotesRequest) (*desc.GetAllNotesResponse, error) {
	var notes []models.Note

	// results with pagination
	result := initializers.DB.Limit(int(req.Limit)).Offset(int(req.Offset)).Find(&notes)

	if result.Error != nil {
		return nil, result.Error
	}

	notesAmount := int32(len(notes))
	// format models.Note to proto.Note
	protoNotes := make([]*desc.Note, notesAmount)
	for _, n := range notes {
		protoNotes = append(protoNotes, &desc.Note{
			Id:      n.ID,
			Title:   n.Title,
			Content: n.Content,
		})
	}

	return &desc.GetAllNotesResponse{
		Notes:      protoNotes,
		TotalCount: notesAmount,
	}, nil
}

func (s *Server) GetNote(ctx context.Context, req *desc.GetNoteRequest) (*desc.Note, error) {
	var note models.Note

	result := initializers.DB.First(&note, req.GetId())
	if result.Error != nil {
		return nil, result.Error
	}

	return &desc.Note{
		Id:      note.ID,
		Title:   note.Title,
		Content: note.Content,
	}, nil
}

func (s *Server) UpdateNote(ctx context.Context, req *desc.UpdateNoteRequest) (*desc.Note, error) {
	oldNote, err := s.GetNote(ctx, &desc.GetNoteRequest{Id: req.GetId()})
	if err != nil {
		return nil, err
	}

	// check if something was updated
	if oldNote.Title == req.GetNewTitle() && oldNote.Content == req.GetNewContent() {
		return nil, status.Errorf(codes.InvalidArgument, "nothing to update")
	}

	// update
	newNote := &models.Note{
		ID:      req.GetId(),
		Title:   req.GetNewTitle(),
		Content: req.GetNewContent(),
	}
	initializers.DB.Model(oldNote).Updates(newNote)

	return &desc.Note{
		Id:      newNote.ID,
		Title:   newNote.Title,
		Content: newNote.Content,
	}, nil
}

func (s *Server) DeleteNote(ctx context.Context, req *desc.DeleteNoteRequest) (*emptypb.Empty, error) {
	result := initializers.DB.Delete(&models.Note{}, req.GetId())

	if result.Error != nil {
		return nil, result.Error
	}

	return &emptypb.Empty{}, nil
}
