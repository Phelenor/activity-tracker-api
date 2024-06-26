package activity

func (activity Activity) ToDbActivity(id string, userId string) DbActivity {
	return DbActivity{
		Id:                        id,
		GroupActivityId:           activity.GroupActivityId,
		UserId:                    userId,
		ActivityType:              activity.ActivityType,
		StartTimestamp:            activity.StartTimestamp,
		EndTimestamp:              activity.EndTimestamp,
		DistanceMeters:            activity.DistanceMeters,
		DurationSeconds:           activity.DurationSeconds,
		AvgSpeedKmh:               activity.AvgSpeedKmh,
		MaxSpeedKmh:               activity.MaxSpeedKmh,
		AvgHeartRate:              activity.AvgHeartRate,
		MaxHeartRate:              activity.MaxHeartRate,
		Calories:                  activity.Calories,
		Elevation:                 activity.Elevation,
		Weather:                   activity.Weather,
		HeartRateZoneDistribution: activity.HeartRateZoneDistribution,
		Goals:                     activity.Goals,
		IsGymActivity:             activity.IsGymActivity,
	}
}

func (dbActivity DbActivity) ToActivity(imageUrl string) Activity {
	return Activity{
		Id:                        dbActivity.Id,
		GroupActivityId:           dbActivity.GroupActivityId,
		ActivityType:              dbActivity.ActivityType,
		StartTimestamp:            dbActivity.StartTimestamp,
		EndTimestamp:              dbActivity.EndTimestamp,
		DistanceMeters:            dbActivity.DistanceMeters,
		DurationSeconds:           dbActivity.DurationSeconds,
		AvgSpeedKmh:               dbActivity.AvgSpeedKmh,
		MaxSpeedKmh:               dbActivity.MaxSpeedKmh,
		AvgHeartRate:              dbActivity.AvgHeartRate,
		MaxHeartRate:              dbActivity.MaxHeartRate,
		Calories:                  dbActivity.Calories,
		Elevation:                 dbActivity.Elevation,
		Weather:                   dbActivity.Weather,
		HeartRateZoneDistribution: dbActivity.HeartRateZoneDistribution,
		Goals:                     dbActivity.Goals,
		IsGymActivity:             dbActivity.IsGymActivity,
		ImageUrl:                  imageUrl,
	}
}
