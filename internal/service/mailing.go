package service

import (
	"context"
	"fmt"
)

func (o *botService) SendingMailingList() {
	for {
		textID := <-o.mailingChan
		text, err := o.repo.Text.GetText(context.Background(), textID)
		if err != nil {
			fmt.Println(fmt.Sprintf("Error getting text: %s", err))
			o.mailingResultChan <- false
			return
		}
		users, err := o.repo.User.GetUsers(context.Background())
		if err != nil {
			fmt.Println(fmt.Sprintf("Error getting users: %s", err))
			o.mailingResultChan <- false
			return
		}
		for _, user := range users {
			if user.IsDeleted {
				continue
			}
			o.sendMessageToChat(user.TelegramID, text.Text)
		}
		o.mailingResultChan <- true
	}
}
