package ZapLogWrapper

import (
	"ZapLogWrapper/constants"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

func InitLogger() *zap.Logger {

	logConfig := zap.NewProductionConfig()

	// CHANGE TO LOG FILE PATH ONLY
	logConfig.OutputPaths = []string{"/Users/subharajbhowmik/go/src/NISAuthenticationService/dummyLog.log", "stderr"}

	logConfig.EncoderConfig.FunctionKey = "func"
	logConfig.EncoderConfig.TimeKey = "time"

	logConfig.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.UTC().Format("2006-01-02T15:04:05Z0700")) // Change to time.Stamp for: Jun 24 11:37:42
	}

	logConfig.DisableStacktrace = true
	logger, _ := logConfig.Build()

	return logger
}

func GenerateLoggerFields(errSrc string, errCode string, subCode string) []zap.Field {
	zapFields :=  []zap.Field{
		{
			Key:    constants.ERROR_SOURCE,
			String: errSrc,
			Type:   zapcore.StringType,
		},
		{
			Key:    constants.ERROR_CODE,
			String: errCode,
			Type:   zapcore.StringType,
		},
	}

	if subCode != "" {
		zapFields = append(zapFields, zap.Field{
			Key:    constants.ERROR_SUB_CODE,
			String: subCode,
			Type:   zapcore.StringType,
		})
	}

	return zapFields
}

func GenerateZapLog(errSrc string, errCode string, errSubCode string, errCause error, addedFields ...interface{}) (string, []zap.Field) {

	zapFields := GenerateLoggerFields(errSrc, errCode, errSubCode)

	switch errSrc {

	case constants.ERR_SRC_JWT:
		switch errCode {
		case constants.ERR_JWT_TOKEN_GENERATION:
			return fmt.Sprintf("Error caught while generating token: %v", errCause), zapFields
		case constants.ERR_INVALID_JWT:
			return fmt.Sprintf("Invalid jwt detected: %v", errCause), zapFields
		case constants.ERR_JWT_PARSE:
			return "Couldn't parse claims", zapFields
		case constants.ERR_EXPIRED_TOKEN:
			return "JWT is expired", zapFields
		case constants.ERR_MISSING_INFO_USER_ID:
			return "Empty id detected", zapFields
		case constants.ERR_MISSING_JWT_SECRET_KEY:
			return "Missing JWT secret key", zapFields
		case constants.ERR_JWT_VALIDATION:
			return "Invalid token detected", zapFields
		default:
			return "JWT Error detected", zapFields
		}

	case constants.ERR_SRC_REQUEST:
		switch errCode {
		case constants.ERR_MISSING_INFO_DEVICE_ID,
			constants.ERR_REQUEST_UNKNOWN_APP,
			constants.ERR_REQUEST_INVALID_PROVIDER,
			constants.ERR_MISSING_INFO_ACCESS_TOKEN,
			constants.ERR_MISSING_INFO_NIS_TOKEN,
			constants.ERR_MISSING_PUID,
			constants.ERR_MISSING_INFO_USER_ID:
			return "Bad request with missing/invalid details", zapFields
		case constants.ERR_PHONE_ALREADY_EXISTS:
			if len(addedFields) == 2 {
				return fmt.Sprintf("Requester %s phone number %s already exists for user-id: %s",
					addedFields[0], errSubCode, addedFields[1]), zapFields
			}
			return "Duplicate phone number detected", zapFields
		default:
			return "Bad Request", zapFields
		}

	case constants.ERR_SRC_DATA_SOURCE:
		switch errCode {
		case constants.ERR_MISSING_EMAIL_AND_PUID:
			return "Data missing for request", zapFields
		default:
			return "Missing data from provider data source", zapFields
		}

	case constants.ERR_SRC_GOOGLE, constants.ERR_SRC_FACEBOOK:
		switch errCode {
		case constants.ERR_EXT_API_REQUEST:
			return fmt.Sprintf("Error caught while making external API request: %v", errCause), zapFields
		case constants.ERR_EXT_API_RESPONSE_FAILURE:
			if len(addedFields) == 1 {
				return fmt.Sprintf("External API failure. Response: %v error: %v", addedFields[0], errCause),zapFields
			}
			return "External API non-successful response", zapFields
		case constants.ERR_EXT_API_RESPONSE_PARSE:
			return fmt.Sprintf("Error caught while parsing external API response: %v", errCause), zapFields
		case constants.ERR_EXT_API_RESPONSE_BODY_CLOSE:
			return fmt.Sprintf("Error caught while closing response body: %v", errCause) ,zapFields
		case constants.ERR_MISSING_CLIENT:
			return fmt.Sprintf("Error getting provider client: %v", errCause), zapFields
		case constants.ERR_AUTH_FAIL:
			return fmt.Sprintf("Mismatched token detected: %v", errCause), zapFields
		default:
			return "Error detected in user auth flow via external provider", zapFields
		}

	case constants.ERR_SRC_FIREBASE:
		switch errCode {
		case constants.ERR_MISSING_CONFIG:
			return "Missing config", zapFields
		case constants.ERR_MISSING_CLIENT:
			return fmt.Sprintf("Error getting auth client: %v", errCause), zapFields
		case constants.ERR_AUTH_FAIL:
			if len(addedFields) == 1 {
				if errCause == nil {
					return fmt.Sprintf("Mismatch fetched token ID: %s", addedFields[0]), zapFields
				}
				return fmt.Sprintf("Error verifying ID token %v, err: %v", addedFields[0], errCause), zapFields
			}
			return fmt.Sprintf("Empty token detected: %v", errCause), zapFields
		case constants.ERR_EXPIRED_TOKEN:
			return "Token has expired", zapFields
		case constants.ERR_MISSING_RECORD:
			return fmt.Sprintf("Error getting user record for user id: %s, error: %v", errSubCode, errCause), zapFields
		case constants.ERR_CLIENT_BUILD_FAILURE:
			return fmt.Sprintf("Failed to create firebase client: %v", errCause), zapFields
		default:
			return "Error detected in user auth flow via firebase", zapFields
		}

	case constants.ERR_SRC_MONGO:
		switch errCode {
		case constants.ERR_MGO_READ_INTERNAL:
			return fmt.Sprintf("Error Caught while reading from mongo: %v", errCause), zapFields
		case constants.ERR_MGO_BULK_WRITE_INTERNAL:
			return fmt.Sprintf("Error while dumping to DB: %v", errCause), zapFields
		case constants.ERR_MGO_WRITE_INTERNAL:
			return fmt.Sprintf("Error while writing to DB: %v", errCause), zapFields
		default:
			return fmt.Sprintf("Error caught from mongo: %v", errCause), zapFields
		}
	}

	return "Unmatched error detected", zapFields
}

