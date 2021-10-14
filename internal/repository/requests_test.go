package repository

//func TestRequestsRepository_InsertRequest(t *testing.T) {
//	db, mock := NewMock(t)
//	defer db.Close()
//
//	type input struct {
//		userID       string
//		sourceID     string
//		sourceFormat string
//		targetFormat string
//		ratio        int
//	}
//
//	requestsRepo := NewRequestsRepository(db)
//	const query = `insert into converter.requests (.+)`
//
//	testTable := []struct {
//		name              string
//		args              input
//		expectedRequestID string
//		mockBehavior      func(input)
//		isErrorExpected   bool
//	}{
//		{
//			name: "Ok",
//			args: input{
//				userID:       "1",
//				sourceID:     "1",
//				sourceFormat: "jpg",
//				targetFormat: "png",
//				ratio:        95,
//			},
//			mockBehavior: func(args input) {
//				rows := sqlmock.NewRows([]string{"id"}).AddRow("1")
//				mock.ExpectQuery(query).
//					WithArgs(args.userID, args.sourceFormat, args.targetFormat, args.ratio).
//					WillReturnRows(rows)
//			},
//			isErrorExpected: false,
//		},
//		//{ // todo
//		//	name: "Database error",
//		//	args: input{
//		//		userID:       "1",
//		//		sourceID:     "1",
//		//		sourceFormat: "jpeg",
//		//		targetFormat: "png",
//		//		ratio:        80,
//		//	},
//		//	mockBehavior: func(args input) {
//		//		rows := sqlmock.NewRows([]string{"id"})
//		//		mock.ExpectQuery(query).
//		//			WithArgs(args.userID, args.sourceFormat, args.targetFormat, args.ratio).
//		//			WillReturnRows(rows)
//		//	},
//		//	isErrorExpected: true,
//		//},
//	}
//
//	for _, tc := range testTable {
//		t.Run(tc.name, func(t *testing.T) {
//			tc.mockBehavior(tc.args)
//			result, err := requestsRepo.InsertRequest(context.TODO(),
//				tc.args.userID, tc.args.sourceID, tc.args.sourceFormat, tc.args.targetFormat, tc.args.ratio)
//			if tc.isErrorExpected {
//				assert.Error(t, err)
//			} else {
//				assert.NoError(t, err)
//				assert.Equal(t, tc.expectedRequestID, result)
//			}
//		})
//	}
//}

//func TestRequestsRepository_GetRequestsByUserID(t *testing.T) {
//	db, mock := NewMock(t)
//	defer db.Close()
//
//	const query = `select (.+) from converter.requests`
//	testTime := time.Now()
//
//	requestsRepo := NewRequestsRepository(db)
//
//	testTable := []struct {
//		name             string
//		userID           string
//		expectedRequests []ConversionRequest
//		mockBehavior     func(string)
//		isErrorExpected  bool
//	}{
//		{
//			name:   "Ok",
//			userID: "1",
//			expectedRequests: []ConversionRequest{
//				{"1", "1", "11", "22", "jpg", "png",
//					95, testTime, testTime, "done"},
//				{"2", "1", "12", "23", "jpg", "png",
//					80, testTime, testTime, "processing"},
//			},
//			mockBehavior: func(userID string) {
//				rows := sqlmock.NewRows([]string{"id", "user_id", "source_id", "target_id, source_format", "target_format",
//					"ratio", "created", "updated", "status"}).
//					//AddRow(mock2.Anything, mock2.Anything, mock2.Anything, mock2.Anything, mock2.Anything, mock2.Anything, mock2.Anything, mock2.Anything, mock2.Anything, mock2.Anything)
//					AddRow("1", "1", "11", "22", "jpg", "png", 95, testTime, testTime, "done")
//				mock.ExpectQuery(query).
//					WithArgs(userID).
//					WillReturnRows(rows)
//			},
//			isErrorExpected: false,
//		},
//	}
//
//	for _, tc := range testTable {
//		t.Run(tc.name, func(t *testing.T) {
//			tc.mockBehavior(tc.userID)
//
//			result, err := requestsRepo.GetUsersRequests(context.TODO(), tc.userID)
//			if tc.isErrorExpected {
//				assert.Error(t, err)
//			} else {
//				assert.NoError(t, err)
//				assert.Equal(t, tc.expectedRequests, result)
//			}
//		})
//	}
//}
