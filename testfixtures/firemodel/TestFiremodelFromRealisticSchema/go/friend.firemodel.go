// DO NOT EDIT - Code generated by firemodel (dev).

package firemodel

import "time"

type Friend struct {
	Username      string    `firestore:"username,omitempty"`
	DisplayName   string    `firestore:"displayName,omitempty"`
	Avatar        *Avatar   `firestore:"avatar,omitempty"`
	FriendsSinice time.Time `firestore:"friendsSinice,omitempty"`
}