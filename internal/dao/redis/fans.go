package redis

func DoFollow(myId, publisherId string) error {
	pipeline := rdb.TxPipeline()
	rdb.Incr(MY_FOLLOWS_COUNTS + ":" + myId)
	rdb.Incr(MY_FANS_COUNTS + ":" + publisherId)
	rdb.Set(FANS_TO_VLOGGER_FOLLOW+":"+myId+":"+publisherId, "1", 0)
	_, err := pipeline.Exec()
	return err
}

func UnFollow(myId, publisherId string) error {
	pipeline := rdb.TxPipeline()
	rdb.Decr(MY_FOLLOWS_COUNTS + ":" + myId)
	rdb.Decr(MY_FANS_COUNTS + ":" + publisherId)
	rdb.Del(FANS_TO_VLOGGER_FOLLOW + ":" + myId + ":" + publisherId)
	_, err := pipeline.Exec()
	return err
}

func GetFansRelations(myId, publisherId string) string {
	result, _ := rdb.Get(FANS_TO_VLOGGER_FOLLOW + ":" + myId + ":" + publisherId).Result()
	return result
}
