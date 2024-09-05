package middlewares

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockdb "github.com/Luckny/space-it/db/mock"
	db "github.com/Luckny/space-it/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAuditLogger(t *testing.T) {
	reqLog := db.RequestLog{
		ID:        uuid.New(),
		Method:    "GET",
		Path:      "/somepath",
		CreatedAt: pgtype.Timestamp{Time: time.Now()},
	}

	resLog := db.ResponseLog{
		ID:        reqLog.ID,
		Status:    http.StatusOK,
		CreatedAt: pgtype.Timestamp{Time: time.Now()},
	}

	user, unHashedPassword := mockdb.RandomUser(t)

	testCases := []struct {
		name          string
		path          string
		setHeader     bool
		username      string
		password      string
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "unauthenticated request -> ok",
			path: reqLog.Path,
			buildStubs: func(store *mockdb.MockStore) {

				reqLogArg := db.CreateUnauthenticatedRequestLogParams{
					Path:   reqLog.Path,
					Method: reqLog.Method,
				}

				resLogArg := db.CreateResponseLogParams{
					ID:     resLog.ID,
					Status: resLog.Status,
				}
				store.EXPECT().
					CreateUnauthenticatedRequestLog(gomock.Any(), mockdb.EqUnauthenticatedLogParam(reqLogArg)).
					Times(1).
					Return(reqLog, nil)
				store.EXPECT().
					CreateResponseLog(gomock.Any(), mockdb.EqResponseLogParam(resLogArg)).
					Times(1).
					Return(resLog, nil)
			},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusOK)
			},
		},

		{
			name:      "authenticated request -> ok",
			path:      reqLog.Path,
			setHeader: true,
			username:  user.Email,
			password:  unHashedPassword,
			buildStubs: func(store *mockdb.MockStore) {

				reqLogArg := db.CreateAuthenticatedRequestLogParams{
					Path:   reqLog.Path,
					Method: reqLog.Method,
					UserID: user.ID,
				}

				resLogArg := db.CreateResponseLogParams{
					ID:     resLog.ID,
					Status: resLog.Status,
				}

				store.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Eq(user.Email)).
					Times(1).
					Return(user, nil)

				store.EXPECT().
					CreateAuthenticatedRequestLog(gomock.Any(), mockdb.EqAuthenticatedLogParam(reqLogArg)).
					Times(1).
					Return(reqLog, nil)
				store.EXPECT().
					CreateResponseLog(gomock.Any(), mockdb.EqResponseLogParam(resLogArg)).
					Times(1).
					Return(resLog, nil)
			},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusOK)
			},
		},

		{
			name: "request log internal error",
			path: reqLog.Path,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUnauthenticatedRequestLog(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.RequestLog{}, sql.ErrConnDone)
				store.EXPECT().
					CreateResponseLog(gomock.Any(), gomock.Any()).
					Times(0)
			},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusInternalServerError)
			},
		},

		{
			name: "response log internal error",
			path: reqLog.Path,
			buildStubs: func(store *mockdb.MockStore) {
				reqLogArg := db.CreateUnauthenticatedRequestLogParams{
					Path:   reqLog.Path,
					Method: reqLog.Method,
				}

				resLogArg := db.CreateResponseLogParams{
					ID:     resLog.ID,
					Status: resLog.Status,
				}

				store.EXPECT().
					CreateUnauthenticatedRequestLog(gomock.Any(), mockdb.EqUnauthenticatedLogParam(reqLogArg)).
					Times(1).
					Return(reqLog, nil)
				store.EXPECT().
					CreateResponseLog(gomock.Any(), mockdb.EqResponseLogParam(resLogArg)).
					Times(1).
					Return(db.ResponseLog{}, sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusOK)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			router := gin.Default()

			router.Use(Authenticate(store))
			router.Use(AuditLogger(store))

			router.GET(tc.path, func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, nil)
			})

			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, tc.path, nil)
			require.NoError(t, err)

			if tc.setHeader {
				request.SetBasicAuth(tc.username, tc.password)
			}

			router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)

		})
	}
}
