package tts

import (
	"fmt"
	"syscall"
	"unsafe"
)

const (
	MSP_SUCCESS = 0
)

func unsafeString(p string) unsafe.Pointer {
	return unsafe.Pointer(&[]byte(p + "\x00")[0])
}

func uintptrToBytes(ptr uintptr, size int) []byte {
	if ptr != 0 {
		us := make([]byte, 0, size)
		p := ptr
		for i := 0; i < size; i++ {
			u := *(*byte)(unsafe.Pointer(p + uintptr(i)))
			us = append(us, u)
		}

		return us
	}

	return nil
}

func uintptrToString(cstr uintptr) string {
	if cstr != 0 {
		us := make([]byte, 0, 256*1024)
		for p := cstr; ; p += 1 {
			u := *(*byte)(unsafe.Pointer(p))
			if u == 0 {
				return string(us)
			}
			us = append(us, u)
		}
	}
	return ""
}

var (
	xunfei                  = syscall.NewLazyDLL("msc_x64.dll")
	procMSPDownloadData     = xunfei.NewProc("MSPDownloadData")
	procMSPGetParam         = xunfei.NewProc("MSPGetParam")
	procMSPGetVersion       = xunfei.NewProc("MSPGetVersion")
	procMSPLogin            = xunfei.NewProc("MSPLogin")
	procMSPLogout           = xunfei.NewProc("MSPLogout")
	procMSPNlpSchCancel     = xunfei.NewProc("MSPNlpSchCancel")
	procMSPNlpSearch        = xunfei.NewProc("MSPNlpSearch")
	procMSPRegisterNotify   = xunfei.NewProc("MSPRegisterNotify")
	procMSPSearch           = xunfei.NewProc("MSPSearch")
	procMSPSetParam         = xunfei.NewProc("MSPSetParam")
	procMSPUploadData       = xunfei.NewProc("MSPUploadData")
	procQISEAudioWrite      = xunfei.NewProc("QISEAudioWrite")
	procQISEGetParam        = xunfei.NewProc("QISEGetParam")
	procQISEGetResult       = xunfei.NewProc("QISEGetResult")
	procQISEResultInfo      = xunfei.NewProc("QISEResultInfo")
	procQISESessionBegin    = xunfei.NewProc("QISESessionBegin")
	procQISESessionEnd      = xunfei.NewProc("QISESessionEnd")
	procQISETextPut         = xunfei.NewProc("QISETextPut")
	procQISRAudioWrite      = xunfei.NewProc("QISRAudioWrite")
	procQISRBuildGrammar    = xunfei.NewProc("QISRBuildGrammar")
	procQISRGetBinaryResult = xunfei.NewProc("QISRGetBinaryResult")
	procQISRGetParam        = xunfei.NewProc("QISRGetParam")
	procQISRGetResult       = xunfei.NewProc("QISRGetResult")
	procQISRSessionBegin    = xunfei.NewProc("QISRSessionBegin")
	procQISRSessionEnd      = xunfei.NewProc("QISRSessionEnd")
	procQISRSetParam        = xunfei.NewProc("QISRSetParam")
	procQISRUpdateLexicon   = xunfei.NewProc("QISRUpdateLexicon")
	procQIVWAudioWrite      = xunfei.NewProc("QIVWAudioWrite")
	procQIVWRegisterNotify  = xunfei.NewProc("QIVWRegisterNotify")
	procQIVWResMerge        = xunfei.NewProc("QIVWResMerge")
	procQIVWSessionBegin    = xunfei.NewProc("QIVWSessionBegin")
	procQIVWSessionEnd      = xunfei.NewProc("QIVWSessionEnd")
	procQTTSAudioGet        = xunfei.NewProc("QTTSAudioGet")
	procQTTSAudioInfo       = xunfei.NewProc("QTTSAudioInfo")
	procQTTSGetParam        = xunfei.NewProc("QTTSGetParam")
	procQTTSSessionBegin    = xunfei.NewProc("QTTSSessionBegin")
	procQTTSSessionEnd      = xunfei.NewProc("QTTSSessionEnd")
	procQTTSSetParam        = xunfei.NewProc("QTTSSetParam")
	procQTTSTextPut         = xunfei.NewProc("QTTSTextPut")
)

type NLPSearchCB func()
type MSP_STATUS_NTF_HANDLER func()
type GrammarCallBack func()
type LexiconCallBack func()
type MsgProcCallBack func()

// const void* MSPAPI MSPDownloadData(const char* params, unsigned int* dataLen, int* errorCode);
func MSPDownloadData(params string) ([]byte, error) {
	var dataLen int
	var errorCode int

	r1, _, _ := procMSPDownloadData.Call(
		uintptr(unsafeString(params)),
		uintptr(unsafe.Pointer(&dataLen)),
		uintptr(unsafe.Pointer(&errorCode)))
	if errorCode != MSP_SUCCESS {
		return nil, fmt.Errorf("MSPDownloadData failed, error %d", errorCode)
	}

	return uintptrToBytes(r1, dataLen), nil
}

// int MSPAPI MSPGetParam( const char *paramName, char *paramValue, unsigned int *valueLen );
func MSPGetParam(paramName string) (string, error) {
	var valueLen uint
	var pValue = make([]byte, 32)

	r1, _, _ := procMSPGetParam.Call(
		uintptr(unsafeString(paramName)),
		uintptr(unsafe.Pointer(&pValue[0])),
		uintptr(unsafe.Pointer(&valueLen)))
	if r1 != MSP_SUCCESS {
		return "", fmt.Errorf("MSPGetParam failed, error %d", r1)
	}

	return string(pValue), nil
}

// const char* MSPAPI MSPGetVersion(const char *verName, int *errorCode);
func MSPGetVersion(verName string) (string, error) {
	var errorCode int

	r1, _, _ := procMSPGetVersion.Call(
		uintptr(unsafeString(verName)),
		uintptr(unsafe.Pointer(&errorCode)))
	if errorCode != MSP_SUCCESS {
		return "", fmt.Errorf("MSPGetVersion failed, error %d", errorCode)
	}

	return uintptrToString(r1), nil
}

// int MSPAPI MSPLogin(const char* usr, const char* pwd, const char* params);
func MSPLogin(params string) error {
	r1, _, _ := procMSPLogin.Call(
		uintptr(0),
		uintptr(0),
		uintptr(unsafeString(params)))
	if r1 != MSP_SUCCESS {
		return fmt.Errorf("MSPLogin failed, error %d", int(r1))
	}

	return nil
}

// int MSPAPI MSPLogout();
func MSPLogout() error {

	r1, _, _ := procMSPLogout.Call()
	if r1 != MSP_SUCCESS {
		return fmt.Errorf("MSPLogout failed, error %d", r1)
	}

	return nil
}

// int MSPAPI MSPNlpSchCancel(const char *sessionID, const char *hints);
func MSPNlpSchCancel(sessionID string, hints string) error {

	r1, _, _ := procMSPNlpSchCancel.Call(
		uintptr(unsafeString(sessionID)),
		uintptr(unsafeString(hints)))
	if r1 != MSP_SUCCESS {
		return fmt.Errorf("MSPNlpSchCancel failed, error %d", int(r1))
	}

	return nil
}

// const char* MSPAPI MSPNlpSearch(const char* params, const char* text, unsigned int textLen, int *errorCode, NLPSearchCB callback, void *userData);
// const char* MSPAPI MSPNlpSearch(const char* params,
// 								   const char* text,
// 								   unsigned int textLen,
// 								   int *errorCode,
// 								   NLPSearchCB callback,
// 								   void *userData);
func MSPNlpSearch(params string, text string, callback NLPSearchCB) (string, []byte, error) {
	var errorCode int
	var textLen = len(text)
	var userData = make([]byte, 256*1024)

	r1, _, _ := procMSPNlpSearch.Call(
		uintptr(unsafeString(params)),
		uintptr(unsafeString(text)),
		uintptr(textLen),
		uintptr(unsafe.Pointer(&errorCode)),
		uintptr(syscall.NewCallback(callback)),
		uintptr(unsafe.Pointer(&userData[0])))
	if errorCode != MSP_SUCCESS {
		return "", nil, fmt.Errorf("MSPNlpSearch failed, error %d", int(errorCode))
	}

	return uintptrToString(r1), userData, nil
}

// int MSPAPI MSPRegisterNotify( msp_status_ntf_handler statusCb, void *userData );
func MSPRegisterNotify(statusCb MSP_STATUS_NTF_HANDLER) ([]byte, error) {
	var userData = make([]byte, 256*1024)

	r1, _, _ := procMSPRegisterNotify.Call(
		uintptr(syscall.NewCallback(statusCb)),
		uintptr(unsafe.Pointer(&userData[0])))
	if r1 != MSP_SUCCESS {
		return nil, fmt.Errorf("MSPRegisterNotify failed, error %d", int(r1))
	}

	return userData, nil
}

// const char* MSPAPI MSPSearch(const char* params, const char* text, unsigned int* dataLen, int* errorCode);
func MSPSearch(params string, text string) (string, error) {
	var dataLen uint
	var errorCode int

	r1, _, _ := procMSPSearch.Call(
		uintptr(unsafeString(params)),
		uintptr(unsafeString(text)),
		uintptr(unsafe.Pointer(&dataLen)),
		uintptr(unsafe.Pointer(&errorCode)))
	if errorCode != MSP_SUCCESS {
		return "", fmt.Errorf("MSPSearch failed, error %d", int(errorCode))
	}

	return uintptrToString(r1), nil
}

// int MSPAPI MSPSetParam( const char* paramName, const char* paramValue );
func MSPSetParam(paramName string, paramValue string) error {

	r1, _, _ := procMSPSetParam.Call(
		uintptr(unsafeString(paramName)),
		uintptr(unsafeString(paramValue)))
	if r1 != MSP_SUCCESS {
		return fmt.Errorf("MSPSetParam failed, error %d", int(r1))
	}

	return nil
}

// const char* MSPAPI MSPUploadData(const char* dataName, void* data, unsigned int dataLen, const char* params, int* errorCode);
func MSPUploadData(dataName string, data []byte, params string) (string, error) {
	var errorCode int

	r1, _, _ := procMSPUploadData.Call(
		uintptr(unsafeString(dataName)),
		uintptr(unsafe.Pointer(&data[0])),
		uintptr(len(data)),
		uintptr(unsafeString(params)),
		uintptr(unsafe.Pointer(&errorCode)))
	if errorCode != MSP_SUCCESS {
		return "", fmt.Errorf("MSPUploadData failed, error %d", int(errorCode))
	}

	return uintptrToString(r1), nil
}

// int MSPAPI QISEAudioWrite(const char* sessionID, const void* waveData, unsigned int waveLen, int audioStatus, int *epStatus, int *Status);
func QISEAudioWrite(sessionID string, waveData []byte, audioStatus int) (int, int, error) {
	var epStatus int
	var Status int
	var waveLen = len(waveData)

	r1, _, _ := procQISEAudioWrite.Call(
		uintptr(unsafeString(sessionID)),
		uintptr(unsafe.Pointer(&waveData[0])),
		uintptr(waveLen),
		uintptr(audioStatus),
		uintptr(unsafe.Pointer(&epStatus)),
		uintptr(unsafe.Pointer(&Status)))
	if r1 != MSP_SUCCESS {
		return 0, 0, fmt.Errorf("QISEAudioWrite failed, error %d", int(r1))
	}

	return epStatus, Status, nil
}

// int MSPAPI QISEGetParam(const char* sessionID, const char* paramName, char* paramValue, unsigned int* valueLen);
func QISEGetParam(sessionID string, paramName string) (string, error) {
	var valueLen uint
	var pValue = make([]byte, 64)

	r1, _, _ := procQISEGetParam.Call(
		uintptr(unsafeString(sessionID)),
		uintptr(unsafeString(paramName)),
		uintptr(unsafe.Pointer(&pValue[0])),
		uintptr(unsafe.Pointer(&valueLen)))
	if r1 != MSP_SUCCESS {
		return "", fmt.Errorf("QISEGetParam failed, error %d", int(r1))
	}

	return string(pValue[:valueLen]), nil
}

// const char * MSPAPI QISEGetResult(const char* sessionID, unsigned int* rsltLen, int* rsltStatus, int *errorCode);
func QISEGetResult(sessionID string) (string, int, error) {
	var rsltLen int
	var rsltStatus int
	var errorCode int

	r1, _, _ := procQISEGetResult.Call(
		uintptr(unsafeString(sessionID)),
		uintptr(unsafe.Pointer(&rsltLen)),
		uintptr(unsafe.Pointer(&rsltStatus)),
		uintptr(unsafe.Pointer(&errorCode)))
	if errorCode != MSP_SUCCESS {
		return "", 0, fmt.Errorf("QISEGetResult failed, error %d", errorCode)
	}

	return string(uintptrToBytes(r1, rsltLen)), rsltStatus, nil
}

// const char* MSPAPI QISEResultInfo(const char* sessionID, int *errorCode);
func QISEResultInfo(sessionID string) (string, error) {
	var errorCode int

	r1, _, _ := procQISEResultInfo.Call(
		uintptr(unsafeString(sessionID)),
		uintptr(unsafe.Pointer(&errorCode)))
	if errorCode != MSP_SUCCESS {
		return "", fmt.Errorf("QISEResultInfo failed, error %d", int(errorCode))
	}

	return uintptrToString(r1), nil
}

// const char* MSPAPI QISESessionBegin(const char* params, const char* userModelId, int* errorCode);
func QISESessionBegin(params string, userModelId string) (string, error) {
	var errorCode int

	r1, _, _ := procQISESessionBegin.Call(
		uintptr(unsafeString(params)),
		uintptr(unsafeString(userModelId)),
		uintptr(unsafe.Pointer(&errorCode)))
	if r1 != MSP_SUCCESS {
		return "", fmt.Errorf("QISESessionBegin failed, error %d", int(errorCode))
	}

	return uintptrToString(r1), nil
}

// int MSPAPI QISESessionEnd(const char* sessionID, const char* hints);
func QISESessionEnd(sessionID string, hints string) error {

	r1, _, _ := procQISESessionEnd.Call(
		uintptr(unsafeString(sessionID)),
		uintptr(unsafeString(hints)))
	if r1 != MSP_SUCCESS {
		return fmt.Errorf("QISESessionEnd failed, error %d", int(r1))
	}

	return nil
}

// int MSPAPI QISETextPut(const char* sessionID, const char* textString, unsigned int textLen, const char* params);
func QISETextPut(sessionID string, textString string, params string) error {
	var textLen = len([]byte(textString))

	r1, _, _ := procQISETextPut.Call(
		uintptr(unsafeString(sessionID)),
		uintptr(unsafeString(textString)),
		uintptr(textLen),
		uintptr(unsafeString(params)))
	if r1 != MSP_SUCCESS {
		return fmt.Errorf("QISESessionEnd failed, error %d", int(r1))
	}

	return nil
}

// int MSPAPI QISRAudioWrite(const char* sessionID, const void* waveData, unsigned int waveLen, int audioStatus, int *epStatus, int *recogStatus);
func QISRAudioWrite(sessionID string, waveData []byte, audioStatus int) (int, int, error) {
	var epStatus int
	var recogStatus int
	var data unsafe.Pointer
	var waveLen = len(waveData)

	if waveData == nil {
		data = nil
	} else {
		data = unsafe.Pointer(&waveData[0])
	}

	r1, _, _ := procQISRAudioWrite.Call(
		uintptr(unsafeString(sessionID)),
		uintptr(data),
		uintptr(waveLen),
		uintptr(audioStatus),
		uintptr(unsafe.Pointer(&epStatus)),
		uintptr(unsafe.Pointer(&recogStatus)))
	if r1 != MSP_SUCCESS {
		return 0, 0, fmt.Errorf("QISRAudioWrite failed, error %d", int(r1))
	}

	return epStatus, recogStatus, nil
}

// int MSPAPI QISRBuildGrammar(const char *grammarType, const char *grammarContent, unsigned int grammarLength, const char *params, GrammarCallBack callback, void *userData);
func QISRBuildGrammar(grammarType string, grammarContent string, params string, callback GrammarCallBack) ([]byte, error) {
	var grammarLength = len([]byte(grammarContent))
	var userData = make([]byte, 256*1024)

	r1, _, _ := procQISRBuildGrammar.Call(
		uintptr(unsafeString(grammarType)),
		uintptr(unsafeString(grammarContent)),
		uintptr(grammarLength),
		uintptr(unsafeString(params)),
		syscall.NewCallback(callback),
		uintptr(unsafe.Pointer(&userData[0])))
	if r1 != MSP_SUCCESS {
		return nil, fmt.Errorf("QISRBuildGrammar failed, error %d", r1)
	}

	return userData, nil
}

// const char * MSPAPI QISRGetBinaryResult(const char* sessionID, unsigned int* rsltLen,int* rsltStatus, int waitTime, int *errorCode);
func QISRGetBinaryResult(sessionID string, waitTime int) (string, int, error) {
	var rsltLen int
	var rsltStatus int
	var errorCode int

	r1, _, _ := procQISRGetBinaryResult.Call(
		uintptr(unsafeString(sessionID)),
		uintptr(unsafe.Pointer(&rsltLen)),
		uintptr(unsafe.Pointer(&rsltStatus)),
		uintptr(waitTime),
		uintptr(unsafe.Pointer(&errorCode)))
	if errorCode != MSP_SUCCESS {
		return "", 0, fmt.Errorf("QISRGetBinaryResult failed, error %d", int(errorCode))
	}

	return string(uintptrToBytes(r1, rsltLen)), rsltStatus, nil
}

// int MSPAPI QISRGetParam(const char* sessionID, const char* paramName, char* paramValue, unsigned int* valueLen);
func QISRGetParam(sessionID string, paramName string) (string, error) {
	var valueLen int
	var pValue = make([]byte, 32)

	r1, _, _ := procQISRGetParam.Call(
		uintptr(unsafeString(sessionID)),
		uintptr(unsafeString(paramName)),
		uintptr(unsafe.Pointer(&pValue[0])),
		uintptr(unsafe.Pointer(&valueLen)))
	if r1 != MSP_SUCCESS {
		return "", fmt.Errorf("QISRGetParam failed, error %d", int(r1))
	}

	return string(pValue[:valueLen]), nil
}

// const char * MSPAPI QISRGetResult(const char* sessionID, int* rsltStatus, int waitTime, int *errorCode);
func QISRGetResult(sessionID string, recStat int, waitTime int) (string, int, error) {
	var rsltStatus = recStat
	var errorCode int

	r1, _, _ := procQISRGetResult.Call(
		uintptr(unsafeString(sessionID)),
		uintptr(unsafe.Pointer(&rsltStatus)),
		uintptr(waitTime),
		uintptr(unsafe.Pointer(&errorCode)))
	if errorCode != MSP_SUCCESS {
		return "", 0, fmt.Errorf("QISRGetResult failed, error %d", int(errorCode))
	}

	return uintptrToString(r1), rsltStatus, nil
}

// const char* MSPAPI QISRSessionBegin(const char* grammarList, const char* params, int* errorCode);
func QISRSessionBegin(grammarList string, params string) (string, error) {
	var errorCode int

	r1, _, _ := procQISRSessionBegin.Call(
		uintptr(unsafeString(grammarList)),
		uintptr(unsafeString(params)),
		uintptr(unsafe.Pointer(&errorCode)))
	if errorCode != MSP_SUCCESS {
		return "", fmt.Errorf("QISRSessionBegin failed, error %d", int(errorCode))
	}

	return uintptrToString(r1), nil
}

// int MSPAPI QISRSessionEnd(const char* sessionID, const char* hints);
func QISRSessionEnd(sessionID string, hints string) error {

	r1, _, _ := procQISRSessionEnd.Call(
		uintptr(unsafeString(sessionID)),
		uintptr(0))
	if r1 != MSP_SUCCESS {
		return fmt.Errorf("QISRSessionEnd failed, error %d", int(r1))
	}

	return nil
}

// int MSPAPI QISRSetParam(const char* sessionID, const char* paramName, const char* paramValue);
func QISRSetParam(sessionID string, paramName string, paramValue string) error {

	r1, _, _ := procQISRSetParam.Call(
		uintptr(unsafeString(sessionID)),
		uintptr(unsafeString(paramName)),
		uintptr(unsafeString(paramValue)))
	if r1 != MSP_SUCCESS {
		return fmt.Errorf("QISRSetParam failed, error %d", int(r1))
	}

	return nil
}

// int MSPAPI QISRUpdateLexicon(const char *lexiconName, const char *lexiconContent, unsigned int lexiconLength, const char *params, LexiconCallBack callback, void *userData);
func QISRUpdateLexicon(lexiconName string, lexiconContent string, params string, callback LexiconCallBack) ([]byte, error) {
	var lexiconLength = len([]byte(lexiconContent))
	var userData = make([]byte, 256*1024)

	r1, _, _ := procQISRUpdateLexicon.Call(
		uintptr(unsafeString(lexiconName)),
		uintptr(unsafeString(lexiconContent)),
		uintptr(lexiconLength),
		uintptr(unsafeString(params)),
		uintptr(unsafe.Pointer(&userData[0])))
	if r1 != MSP_SUCCESS {
		return nil, fmt.Errorf("QISRUpdateLexicon failed, error %d", int(r1))
	}

	return userData, nil
}

// int MSPAPI QIVWAudioWrite(const char *sessionID, const void *audioData, unsigned int audioLen, int audioStatus);
func QIVWAudioWrite(sessionID string, audioData []byte, audioStatus int) error {
	var audioLen = len(audioData)

	r1, _, _ := procQIVWAudioWrite.Call(
		uintptr(unsafeString(sessionID)),
		uintptr(unsafe.Pointer(&audioData[0])),
		uintptr(audioLen),
		uintptr(audioStatus))
	if r1 != MSP_SUCCESS {
		return fmt.Errorf("QIVWAudioWrite failed, error %d", int(r1))
	}

	return nil
}

// int MSPAPI QIVWRegisterNotify(const char *sessionID, ivw_ntf_handler msgProcCb, void *userData);
func QIVWRegisterNotify(sessionID string, userData []byte, ivw_ntf_handler MsgProcCallBack) error {

	r1, _, _ := procQIVWRegisterNotify.Call(
		uintptr(unsafeString(sessionID)),
		syscall.NewCallback(ivw_ntf_handler),
		uintptr(unsafe.Pointer(&userData[0])))
	if r1 != MSP_SUCCESS {
		return fmt.Errorf("QIVWRegisterNotify failed, error %d", int(r1))
	}

	return nil
}

// int MSPAPI QIVWResMerge(const char *srcPath, const char *destPath, const char *params);
func QIVWResMerge(srcPath string, destPath string, params string) error {

	r1, _, _ := procQIVWResMerge.Call(
		uintptr(unsafeString(srcPath)),
		uintptr(unsafeString(destPath)),
		uintptr(unsafeString(params)))
	if r1 != MSP_SUCCESS {
		return fmt.Errorf("QIVWResMerge failed, error %d", int(r1))
	}

	return nil
}

// const char* MSPAPI QIVWSessionBegin(const char *grammarList, const char *params, int *errorCode);
func QIVWSessionBegin(grammarList string, params string) (string, error) {
	var errorCode int

	r1, _, _ := procQIVWSessionBegin.Call(
		uintptr(unsafeString(grammarList)),
		uintptr(unsafeString(params)),
		uintptr(unsafe.Pointer(&errorCode)))
	if errorCode != MSP_SUCCESS {
		return "", fmt.Errorf("QIVWSessionBegin failed, error %d", int(errorCode))
	}

	return uintptrToString(r1), nil
}

// int MSPAPI QIVWSessionEnd(const char *sessionID, const char *hints);
func QIVWSessionEnd(sessionID string, hints string) error {

	r1, _, _ := procQIVWSessionEnd.Call(
		uintptr(unsafeString(sessionID)),
		uintptr(unsafeString(hints)))
	if r1 != MSP_SUCCESS {
		return fmt.Errorf("QIVWSessionEnd failed, error %d", int(r1))
	}

	return nil
}

// const void* MSPAPI QTTSAudioGet(const char* sessionID, unsigned int* audioLen, int* synthStatus, int* errorCode);
func QTTSAudioGet(sessionID string) ([]byte, int, error) {
	var audioLen int
	var synthStatus int
	var errorCode int

	r1, _, _ := procQTTSAudioGet.Call(
		uintptr(unsafeString(sessionID)),
		uintptr(unsafe.Pointer(&audioLen)),
		uintptr(unsafe.Pointer(&synthStatus)),
		uintptr(unsafe.Pointer(&errorCode)))
	if errorCode != MSP_SUCCESS {
		return nil, 0, fmt.Errorf("QTTSAudioGet failed, error %d", int(errorCode))
	}

	return uintptrToBytes(r1, audioLen), synthStatus, nil
}

// const char* MSPAPI QTTSAudioInfo(const char* sessionID);
func QTTSAudioInfo(sessionID string) string {

	r1, _, _ := procQTTSAudioInfo.Call(uintptr(unsafeString(sessionID)))
	return uintptrToString(r1)
}

// int MSPAPI QTTSGetParam(const char* sessionID, const char* paramName, char* paramValue, unsigned int* valueLen);
func QTTSGetParam(sessionID string, paramName string) (string, error) {
	var valueLen uint
	var pValue = make([]byte, 32)

	r1, _, _ := procQTTSGetParam.Call(
		uintptr(unsafeString(sessionID)),
		uintptr(unsafeString(paramName)),
		uintptr(unsafe.Pointer(&pValue[0])),
		uintptr(unsafe.Pointer(&valueLen)))
	if r1 != MSP_SUCCESS {
		return "", fmt.Errorf("QTTSGetParam failed, error %d", int(r1))
	}

	return string(pValue[:valueLen]), nil
}

// const char* MSPAPI QTTSSessionBegin(const char* params, int* errorCode);
func QTTSSessionBegin(params string) (string, error) {
	var errorCode int

	r1, _, _ := procQTTSSessionBegin.Call(
		uintptr(unsafeString(params)),
		uintptr(unsafe.Pointer(&errorCode)))
	if errorCode != MSP_SUCCESS {
		return "", fmt.Errorf("QTTSSessionBegin failed, error %d", int(errorCode))
	}

	return uintptrToString(r1), nil
}

// int MSPAPI QTTSSessionEnd(const char* sessionID, const char* hints);
func QTTSSessionEnd(sessionID string, hints string) error {

	r1, _, _ := procQTTSSessionEnd.Call(
		uintptr(unsafeString(sessionID)),
		uintptr(unsafeString(hints)))
	if r1 != MSP_SUCCESS {
		return fmt.Errorf("QTTSSessionEnd failed, error %d", int(r1))
	}

	return nil
}

// int MSPAPI QTTSSetParam(const char *sessionID, const char *paramName, const char *paramValue);
func QTTSSetParam(sessionID string, paramName string, paramValue string) error {

	r1, _, _ := procQTTSSetParam.Call(
		uintptr(unsafeString(sessionID)),
		uintptr(unsafeString(paramName)),
		uintptr(unsafeString(paramValue)))
	if r1 != MSP_SUCCESS {
		return fmt.Errorf("QTTSSetParam failed, error %d", int(r1))
	}

	return nil
}

// int MSPAPI QTTSTextPut(const char* sessionID, const char* textString, unsigned int textLen, const char* params);
func QTTSTextPut(sessionID string, textString string, params string) error {
	var textLen = len([]byte(textString))

	r1, _, _ := procQTTSTextPut.Call(
		uintptr(unsafeString(sessionID)),
		uintptr(unsafeString(textString)),
		uintptr(textLen),
		uintptr(unsafeString(params)))
	if r1 != MSP_SUCCESS {
		return fmt.Errorf("QTTSTextPut failed, error %d", int(r1))
	}

	return nil
}
