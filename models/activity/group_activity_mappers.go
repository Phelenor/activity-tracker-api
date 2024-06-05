package activity

func (db *DbGroupActivity) ToGroupActivity() *GroupActivity {
	return &GroupActivity{
		Id:             db.Id,
		JoinCode:       db.JoinCode,
		UserOwnerId:    db.UserOwnerId,
		UserOwnerName:  db.UserOwnerName,
		ActivityType:   db.ActivityType,
		StartTimestamp: db.StartTimestamp,
		Status:         db.Status,
		JoinedUsers:    db.JoinedUsers,
		ConnectedUsers: db.ConnectedUsers,
		ActiveUsers:    db.ActiveUsers,
		FinishedUsers:  db.FinishedUsers,
	}
}

func (ga *GroupActivity) ToDbGroupActivity() *DbGroupActivity {
	return &DbGroupActivity{
		Id:             ga.Id,
		JoinCode:       ga.JoinCode,
		UserOwnerId:    ga.UserOwnerId,
		UserOwnerName:  ga.UserOwnerName,
		ActivityType:   ga.ActivityType,
		StartTimestamp: ga.StartTimestamp,
		Status:         ga.Status,
		JoinedUsers:    ga.JoinedUsers,
		ConnectedUsers: ga.ConnectedUsers,
		ActiveUsers:    ga.ActiveUsers,
		FinishedUsers:  ga.FinishedUsers,
	}
}
