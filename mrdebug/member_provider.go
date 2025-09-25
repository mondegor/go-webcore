package mrdebug

import (
	"context"

	"github.com/mondegor/go-webcore/mraccess"
)

type (
	// MemberProvider - возвращает сущность с наделёнными ей правами доступа к системе.
	MemberProvider struct {
		member mraccess.Member
	}
)

// NewMemberProvider - создаёт объект MemberProvider.
func NewMemberProvider(member mraccess.Member) *MemberProvider {
	return &MemberProvider{
		member: member,
	}
}

// MemberByToken - comment method.
func (p *MemberProvider) MemberByToken(_ context.Context, _ string) (mraccess.Member, error) {
	return p.member, nil
}
