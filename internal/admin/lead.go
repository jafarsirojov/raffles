package admin

import (
	"context"
	"crm/internal/structs"
	"crm/pkg/errors"
	"fmt"
	"go.uber.org/zap"
	"strings"
	"time"
)

func (s *service) GetLeadList(ctx context.Context, offset, limit int, status string) (list []structs.Lead, err error) {

	if len(strings.TrimSpace(status)) == 0 {
		list, err = s.leadRepo.GetLeadList(ctx, offset, limit)
		if err != nil {
			s.logger.Error("internal.admin.GetLeadList s.leadRepo.GetLeadList", zap.Error(err))
			return nil, err
		}
	} else {
		list, err = s.leadRepo.GetLeadListByStatus(ctx, offset, limit, status)
		if err != nil {
			s.logger.Error("internal.admin.GetLeadList s.leadRepo.GetLeadListByStatus", zap.Error(err))
			return nil, err
		}
	}

	return list, nil
}

func (s *service) GetLeadListXLSX(ctx context.Context, status string) (url string, err error) {

	var list []structs.Lead

	if len(strings.TrimSpace(status)) == 0 {
		list, err = s.leadRepo.GetLeadList(ctx, 0, 1000)
		if err != nil {
			s.logger.Error("internal.admin.GetLeadListXLSX s.leadRepo.GetLeadList", zap.Error(err))
			return url, err
		}
	} else {
		list, err = s.leadRepo.GetLeadListByStatus(ctx, 0, 1000, status)
		if err != nil {
			s.logger.Error("internal.admin.GetLeadListXLSX s.leadRepo.GetLeadListByStatus", zap.Error(err))
			return url, err
		}
	}

	name := time.Now().Unix()

	pathCSV := fmt.Sprintf("lead/%d.csv", name)
	pathXLSX := fmt.Sprintf("lead/%d.xlsx", name)
	dir := "./../../www/harbour.iqomi.ae/"

	err = s.createCSV(dir+pathCSV, list)
	if err != nil {
		s.logger.Error("internal.admin.GetLeadListXLSX s.createCSV", zap.Error(err))
		return "", err
	}

	err = s.generateXLSXFromCSV(dir+pathCSV, dir+pathXLSX)
	if err != nil {
		s.logger.Error("internal.admin.GetLeadListXLSX s.generateXLSXFromCSV", zap.Error(err))
		return "", err
	}

	return URL + "/" + pathXLSX, nil
}

const URL = "harbour.iqomi.ae"

func (s *service) GetLeadAndCommentsByID(ctx context.Context, id int) (lc structs.LeadAndComments, err error) {

	lead, err := s.leadRepo.GetLeadByID(ctx, id)
	if err != nil {
		s.logger.Error("internal.admin.GetLeadAndCommentsByID s.leadRepo.GetLeadByID", zap.Error(err))
		return lc, err
	}

	comments, err := s.commentRepo.GetCommentsByLeadID(ctx, lead.ID)
	if err != nil && err != errors.ErrNotFound {
		s.logger.Error("internal.admin.GetLeadAndCommentsByID s.commentRepo.GetCommentsByLeadID", zap.Error(err))
		return lc, err
	}

	lc.Lead = lead
	lc.Comments = comments

	return lc, nil
}

func (s *service) AddComment(ctx context.Context, comment structs.Comment) error {

	err := s.commentRepo.AddComment(ctx, comment)
	if err != nil {
		s.logger.Error("internal.admin.AddComment s.commentRepo.AddComment", zap.Error(err))
		return err
	}

	return nil
}

func (s *service) UpdateLeadStatus(ctx context.Context, id int, status string) error {

	err := s.leadRepo.UpdateLeadStatus(ctx, id, status)
	if err != nil {
		s.logger.Error("internal.admin.UpdateLeadStatus s.leadRepo.UpdateLeadStatus", zap.Error(err))
		return err
	}

	return nil
}
