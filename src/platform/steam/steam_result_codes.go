/******************************************************************************/
/* steam_result_codes.go                                                      */
/******************************************************************************/
/*                           This file is part of:                            */
/*                                KAIJU ENGINE                                */
/*                          https://kaijuengine.org                           */
/******************************************************************************/
/* MIT License                                                                */
/*                                                                            */
/* Copyright (c) 2023-present Kaiju Engine authors (AUTHORS.md).              */
/* Copyright (c) 2015-present Brent Farris.                                   */
/*                                                                            */
/* May all those that this source may reach be blessed by the LORD and find   */
/* peace and joy in life.                                                     */
/* Everyone who drinks of this water will be thirsty again; but whoever       */
/* drinks of the water that I will give him shall never thirst; John 4:13-14  */
/*                                                                            */
/* Permission is hereby granted, free of charge, to any person obtaining a    */
/* copy of this software and associated documentation files (the "Software"), */
/* to deal in the Software without restriction, including without limitation  */
/* the rights to use, copy, modify, merge, publish, distribute, sublicense,   */
/* and/or sell copies of the Software, and to permit persons to whom the      */
/* Software is furnished to do so, subject to the following conditions:       */
/*                                                                            */
/* The above copyright, blessing, biblical verse, notice and                  */
/* this permission notice shall be included in all copies or                  */
/* substantial portions of the Software.                                      */
/*                                                                            */
/* THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS    */
/* OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF                 */
/* MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.     */
/* IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY       */
/* CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT  */
/* OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE      */
/* OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.                              */
/******************************************************************************/

package steam

type ResultCode int

const (
	ResultCodeNone         = 0 // no result
	ResultCodeOK           = 1 // success
	ResultCodeFail         = 2 // generic failure
	ResultCodeNoConnection = 3 // no/failed network connection
	//	ResultCodeNoConnectionRetry = 4				// OBSOLETE - removed
	ResultCodeInvalidPassword                         = 5  // password/ticket is invalid
	ResultCodeLoggedInElsewhere                       = 6  // same user logged in elsewhere
	ResultCodeInvalidProtocolVer                      = 7  // protocol version is incorrect
	ResultCodeInvalidParam                            = 8  // a parameter is incorrect
	ResultCodeFileNotFound                            = 9  // file was not found
	ResultCodeBusy                                    = 10 // called method busy - action not taken
	ResultCodeInvalidState                            = 11 // called object was in an invalid state
	ResultCodeInvalidName                             = 12 // name is invalid
	ResultCodeInvalidEmail                            = 13 // email is invalid
	ResultCodeDuplicateName                           = 14 // name is not unique
	ResultCodeAccessDenied                            = 15 // access is denied
	ResultCodeTimeout                                 = 16 // operation timed out
	ResultCodeBanned                                  = 17 // VAC2 banned
	ResultCodeAccountNotFound                         = 18 // account not found
	ResultCodeInvalidSteamID                          = 19 // steamID is invalid
	ResultCodeServiceUnavailable                      = 20 // The requested service is currently unavailable
	ResultCodeNotLoggedOn                             = 21 // The user is not logged on
	ResultCodePending                                 = 22 // Request is pending (may be in process, or waiting on third party)
	ResultCodeEncryptionFailure                       = 23 // Encryption or Decryption failed
	ResultCodeInsufficientPrivilege                   = 24 // Insufficient privilege
	ResultCodeLimitExceeded                           = 25 // Too much of a good thing
	ResultCodeRevoked                                 = 26 // Access has been revoked (used for revoked guest passes)
	ResultCodeExpired                                 = 27 // License/Guest pass the user is trying to access is expired
	ResultCodeAlreadyRedeemed                         = 28 // Guest pass has already been redeemed by account, cannot be acked again
	ResultCodeDuplicateRequest                        = 29 // The request is a duplicate and the action has already occurred in the past, ignored this time
	ResultCodeAlreadyOwned                            = 30 // All the games in this guest pass redemption request are already owned by the user
	ResultCodeIPNotFound                              = 31 // IP address not found
	ResultCodePersistFailed                           = 32 // failed to write change to the data store
	ResultCodeLockingFailed                           = 33 // failed to acquire access lock for this operation
	ResultCodeLogonSessionReplaced                    = 34
	ResultCodeConnectFailed                           = 35
	ResultCodeHandshakeFailed                         = 36
	ResultCodeIOFailure                               = 37
	ResultCodeRemoteDisconnect                        = 38
	ResultCodeShoppingCartNotFound                    = 39 // failed to find the shopping cart requested
	ResultCodeBlocked                                 = 40 // a user didn't allow it
	ResultCodeIgnored                                 = 41 // target is ignoring sender
	ResultCodeNoMatch                                 = 42 // nothing matching the request found
	ResultCodeAccountDisabled                         = 43
	ResultCodeServiceReadOnly                         = 44 // this service is not accepting content changes right now
	ResultCodeAccountNotFeatured                      = 45 // account doesn't have value, so this feature isn't available
	ResultCodeAdministratorOK                         = 46 // allowed to take this action, but only because requester is admin
	ResultCodeContentVersion                          = 47 // A Version mismatch in content transmitted within the Steam protocol.
	ResultCodeTryAnotherCM                            = 48 // The current CM can't service the user making a request, user should try another.
	ResultCodePasswordRequiredToKickSession           = 49 // You are already logged in elsewhere, this cached credential login has failed.
	ResultCodeAlreadyLoggedInElsewhere                = 50 // You are already logged in elsewhere, you must wait
	ResultCodeSuspended                               = 51 // Long running operation (content download) suspended/paused
	ResultCodeCancelled                               = 52 // Operation canceled (typically by user: content download)
	ResultCodeDataCorruption                          = 53 // Operation canceled because data is ill formed or unrecoverable
	ResultCodeDiskFull                                = 54 // Operation canceled - not enough disk space.
	ResultCodeRemoteCallFailed                        = 55 // an remote call or IPC call failed
	ResultCodePasswordUnset                           = 56 // Password could not be verified as it's unset server side
	ResultCodeExternalAccountUnlinked                 = 57 // External account (PSN, Facebook...) is not linked to a Steam account
	ResultCodePSNTicketInvalid                        = 58 // PSN ticket was invalid
	ResultCodeExternalAccountAlreadyLinked            = 59 // External account (PSN, Facebook...) is already linked to some other account, must explicitly request to replace/delete the link first
	ResultCodeRemoteFileConflict                      = 60 // The sync cannot resume due to a conflict between the local and remote files
	ResultCodeIllegalPassword                         = 61 // The requested new password is not legal
	ResultCodeSameAsPreviousValue                     = 62 // new value is the same as the old one ( secret question and answer )
	ResultCodeAccountLogonDenied                      = 63 // account login denied due to 2nd factor authentication failure
	ResultCodeCannotUseOldPassword                    = 64 // The requested new password is not legal
	ResultCodeInvalidLoginAuthCode                    = 65 // account login denied due to auth code invalid
	ResultCodeAccountLogonDeniedNoMail                = 66 // account login denied due to 2nd factor auth failure - and no mail has been sent - partner site specific
	ResultCodeHardwareNotCapableOfIPT                 = 67 //
	ResultCodeIPTInitError                            = 68 //
	ResultCodeParentalControlRestricted               = 69 // operation failed due to parental control restrictions for current user
	ResultCodeFacebookQueryError                      = 70 // Facebook query returned an error
	ResultCodeExpiredLoginAuthCode                    = 71 // account login denied due to auth code expired
	ResultCodeIPLoginRestrictionFailed                = 72
	ResultCodeAccountLockedDown                       = 73
	ResultCodeAccountLogonDeniedVerifiedEmailRequired = 74
	ResultCodeNoMatchingURL                           = 75
	ResultCodeBadResponse                             = 76  // parse failure, missing field, etc.
	ResultCodeRequirePasswordReEntry                  = 77  // The user cannot complete the action until they re-enter their password
	ResultCodeValueOutOfRange                         = 78  // the value entered is outside the acceptable range
	ResultCodeUnexpectedError                         = 79  // something happened that we didn't expect to ever happen
	ResultCodeDisabled                                = 80  // The requested service has been configured to be unavailable
	ResultCodeInvalidCEGSubmission                    = 81  // The set of files submitted to the CEG server are not valid !
	ResultCodeRestrictedDevice                        = 82  // The device being used is not allowed to perform this action
	ResultCodeRegionLocked                            = 83  // The action could not be complete because it is region restricted
	ResultCodeRateLimitExceeded                       = 84  // Temporary rate limit exceeded, try again later, different from k_EResultLimitExceeded which may be permanent
	ResultCodeAccountLoginDeniedNeedTwoFactor         = 85  // Need two-factor code to login
	ResultCodeItemDeleted                             = 86  // The thing we're trying to access has been deleted
	ResultCodeAccountLoginDeniedThrottle              = 87  // login attempt failed, try to throttle response to possible attacker
	ResultCodeTwoFactorCodeMismatch                   = 88  // two factor code mismatch
	ResultCodeTwoFactorActivationCodeMismatch         = 89  // activation code for two-factor didn't match
	ResultCodeAccountAssociatedToMultiplePartners     = 90  // account has been associated with multiple partners
	ResultCodeNotModified                             = 91  // data not modified
	ResultCodeNoMobileDevice                          = 92  // the account does not have a mobile device associated with it
	ResultCodeTimeNotSynced                           = 93  // the time presented is out of range or tolerance
	ResultCodeSmsCodeFailed                           = 94  // SMS code failure (no match, none pending, etc.)
	ResultCodeAccountLimitExceeded                    = 95  // Too many accounts access this resource
	ResultCodeAccountActivityLimitExceeded            = 96  // Too many changes to this account
	ResultCodePhoneActivityLimitExceeded              = 97  // Too many changes to this phone
	ResultCodeRefundToWallet                          = 98  // Cannot refund to payment method, must use wallet
	ResultCodeEmailSendFailure                        = 99  // Cannot send an email
	ResultCodeNotSettled                              = 100 // Can't perform operation till payment has settled
	ResultCodeNeedCaptcha                             = 101 // Needs to provide a valid captcha
	ResultCodeGSLTDenied                              = 102 // a game server login token owned by this token's owner has been banned
	ResultCodeGSOwnerDenied                           = 103 // game server owner is denied for other reason (account lock, community ban, vac ban, missing phone)
	ResultCodeInvalidItemType                         = 104 // the type of thing we were requested to act on is invalid
	ResultCodeIPBanned                                = 105 // the ip address has been banned from taking this action
	ResultCodeGSLTExpired                             = 106 // this token has expired from disuse; can be reset for use
	ResultCodeInsufficientFunds                       = 107 // user doesn't have enough wallet funds to complete the action
	ResultCodeTooManyPending                          = 108 // There are too many of this thing pending already
	ResultCodeNoSiteLicensesFound                     = 109 // No site licenses found
	ResultCodeWGNetworkSendExceeded                   = 110 // the WG couldn't send a response because we exceeded max network send size
	ResultCodeAccountNotFriends                       = 111 // the user is not mutually friends
	ResultCodeLimitedUserAccount                      = 112 // the user is limited
	ResultCodeCantRemoveItem                          = 113 // item can't be removed
	ResultCodeAccountDeleted                          = 114 // account has been deleted
	ResultCodeExistingUserCancelledLicense            = 115 // A license for this already exists, but cancelled
	ResultCodeCommunityCooldown                       = 116 // access is denied because of a community cooldown (probably from support profile data resets)
	ResultCodeNoLauncherSpecified                     = 117 // No launcher was specified, but a launcher was needed to choose correct realm for operation.
	ResultCodeMustAgreeToSSA                          = 118 // User must agree to china SSA or global SSA before login
	ResultCodeLauncherMigrated                        = 119 // The specified launcher type is no longer supported; the user should be directed elsewhere
	ResultCodeSteamRealmMismatch                      = 120 // The user's realm does not match the realm of the requested resource
	ResultCodeInvalidSignature                        = 121 // signature check did not match
	ResultCodeParseFailure                            = 122 // Failed to parse input
	ResultCodeNoVerifiedPhone                         = 123 // account does not have a verified phone number
	ResultCodeInsufficientBattery                     = 124 // user device doesn't have enough battery charge currently to complete the action
	ResultCodeChargerRequired                         = 125 // The operation requires a charger to be plugged in, which wasn't present
	ResultCodeCachedCredentialInvalid                 = 126 // Cached credential was invalid - user must reauthenticate
	ResultCodePhoneNumberIsVOIP                       = 127 // The phone number provided is a Voice Over IP number
	ResultCodeNotSupported                            = 128 // The data being accessed is not supported by this API
	ResultCodeFamilySizeLimitExceeded                 = 129 // Reached the maximum size of the family
)
