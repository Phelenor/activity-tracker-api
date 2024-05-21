package activity

func (activity Activity) ToDbActivity(id string, userId string) DbActivity {
	return DbActivity{
		Id:                        id,
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
	}
}

func (dbActivity DbActivity) ToActivity(imageUrl string) Activity {
	return Activity{
		Id:                        dbActivity.Id,
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
		ImageUrl:                  imageUrl,
	}
}
