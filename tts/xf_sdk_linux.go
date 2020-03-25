package tts

/*
#cgo CFLAGS:-g -Wall -I./include
#cgo linux LDFLAGS:-L./libs -lmsc -lrt -ldl -lpthread -lstdc++
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <errno.h>
#include "include/qtts.h"
#include "include/msp_cmn.h"
#include "include/msp_errors.h"
#include "include/msp_types.h"
#ifdef __linux__
#include <unistd.h>
#include <sys/time.h>
#endif

*/
import "C"
import (
	"fmt"
	"unsafe"
)

type NLPSearchCB C.NLPSearchCB
type MSP_STATUS_NTF_HANDLER C.msp_status_ntf_handler

func CString(str string) *C.char {
	return C.CString(str + string(0))
}

// const void* MSPAPI MSPDownloadData(const char* params, unsigned int* dataLen, int* errorCode);
func MSPDownloadData(params string) ([]byte, error) {
	var dataLen C.uint
	var errorCode C.int

	ret := C.MSPDownloadData(CString(params), &dataLen, &errorCode)
	if errorCode != C.MSP_SUCCESS {
		return nil, fmt.Errorf("MSPDownloadData failed, error %d", int(errorCode))
	}

	return C.GoBytes(ret, C.int(dataLen)), nil
}

// int MSPAPI MSPGetParam( const char *paramName, char *paramValue, unsigned int *valueLen );
func MSPGetParam(paramName string) (string, error) {
	var valueLen C.uint
	var paramValue [32]C.char

	ret := C.MSPGetParam(CString(paramName), &paramValue[0], &valueLen)
	if ret != C.MSP_SUCCESS {
		return "", fmt.Errorf("MSPGetParam failed, error %d", int(ret))
	}

	return C.GoStringN(&paramValue[0], C.int(valueLen)), nil
}

// const char* MSPAPI MSPGetVersion(const char *verName, int *errorCode);
func MSPGetVersion(verName string) (string, error) {
	var errorCode C.int

	ret := C.MSPGetVersion(CString(verName), &errorCode)
	if errorCode != C.MSP_SUCCESS {
		return "", fmt.Errorf("MSPGetVersion failed, error %d", int(errorCode))
	}

	return C.GoString(ret), nil
}

// int MSPAPI MSPLogin(const char* usr, const char* pwd, const char* params);
func MSPLogin(params string) error {

	ret := C.MSPLogin(nil, nil, CString(params))
	if ret != C.MSP_SUCCESS {
		return fmt.Errorf("MSPLogin failed, error %d", int(ret))
	}

	return nil
}

// int MSPAPI MSPLogout();
func MSPLogout() error {

	ret := C.MSPLogout()
	if ret != C.MSP_SUCCESS {
		return fmt.Errorf("MSPLogout failed, error %d", int(ret))
	}

	return nil
}

// int MSPAPI MSPNlpSchCancel(const char *sessionID, const char *hints);
func MSPNlpSchCancel(sessionID string, hints string) error {

	ret := C.MSPNlpSchCancel(CString(sessionID), CString(hints))
	if ret != C.MSP_SUCCESS {
		return fmt.Errorf("MSPNlpSchCancel failed, error %d", int(ret))
	}

	return nil
}

// const char* MSPAPI MSPNlpSearch(const char* params, const char* text, unsigned int textLen, int *errorCode, NLPSearchCB callback, void *userData);
func MSPNlpSearch(params string, text string, callback NLPSearchCB) (string, []byte, error) {
	var errorCode C.int
	var txt = CString(text)
	var textLen = C.strlen(txt)
	var userData = make([]byte, 256*1024)

	ret := C.MSPNlpSearch(
		CString(params),
		txt,
		C.uint(textLen),
		&errorCode,
		C.NLPSearchCB(callback),
		unsafe.Pointer(&userData[0]),
	)
	if errorCode != C.MSP_SUCCESS {
		return "", nil, fmt.Errorf("MSPNlpSearch failed, error %d", int(errorCode))
	}

	return C.GoString(ret), userData, nil
}

// int MSPAPI MSPRegisterNotify( msp_status_ntf_handler statusCb, void *userData );
func MSPRegisterNotify(statusCb MSP_STATUS_NTF_HANDLER) ([]byte, error) {
	var userData = make([]byte, 256*1024)
	ret := C.MSPRegisterNotify(statusCb, unsafe.Pointer(&userData[0]))

	if ret != C.MSP_SUCCESS {
		return nil, fmt.Errorf("MSPRegisterNotify failed, error %d", int(ret))
	}

	return userData, nil
}

// const char* MSPAPI MSPSearch(const char* params, const char* text, unsigned int* dataLen, int* errorCode);
func MSPSearch(params string, text string) (string, error) {
	var dataLen C.uint
	var errorCode C.int

	ret := C.MSPSearch(CString(params), CString(text), &dataLen, &errorCode)
	if errorCode != C.MSP_SUCCESS {
		return "", fmt.Errorf("MSPSearch failed, error %d", int(errorCode))
	}

	return C.GoString(ret), nil
}

// int MSPAPI MSPSetParam( const char* paramName, const char* paramValue );
func MSPSetParam(paramName string, paramValue string) error {

	ret := C.MSPSetParam(CString(paramName), CString(paramValue))
	if ret != C.MSP_SUCCESS {
		return fmt.Errorf("MSPSetParam failed, error %d", int(ret))
	}

	return nil
}

// const char* MSPAPI MSPUploadData(const char* dataName, void* data, unsigned int dataLen, const char* params, int* errorCode);
func MSPUploadData(dataName string, data []byte, params string) (string, error) {
	var errorCode C.int

	ret := C.MSPUploadData(CString(dataName), C.CBytes(data), C.uint(len(data)), C.CString(params), &errorCode)
	if errorCode != C.MSP_SUCCESS {
		return "", fmt.Errorf("MSPUploadData failed, error %d", int(errorCode))
	}

	return C.GoString(ret), nil
}

// const void* MSPAPI QTTSAudioGet(const char* sessionID, unsigned int* audioLen, int* synthStatus, int* errorCode);
func QTTSAudioGet(sessionID string) ([]byte, int, error) {
	var audioLen C.uint
	var synthStatus C.int
	var errorCode C.int

	ret := C.QTTSAudioGet(C.CString(sessionID), &audioLen, &synthStatus, &errorCode)
	if errorCode != C.MSP_SUCCESS {
		return nil, 0, fmt.Errorf("QTTSAudioGet failed, error %d", int(errorCode))
	}

	return C.GoBytes(ret, C.int(audioLen)), int(synthStatus), nil
}

// const char* MSPAPI QTTSAudioInfo(const char* sessionID);
func QTTSAudioInfo(sessionID string) string {
	return C.GoString(C.QTTSAudioInfo(C.CString(sessionID)))
}

// int MSPAPI QTTSGetParam(const char* sessionID, const char* paramName, char* paramValue, unsigned int* valueLen);
func QTTSGetParam(sessionID string, paramName string) (string, error) {
	var valueLen C.uint
	var paramValue [32]C.char

	ret := C.QTTSGetParam(C.CString(sessionID), C.CString(paramName), &paramValue[0], &valueLen)
	if ret != C.MSP_SUCCESS {
		return "", fmt.Errorf("QTTSGetParam failed, error %d", int(ret))
	}

	return C.GoStringN(&paramValue[0], C.int(valueLen)), nil
}

// const char* MSPAPI QTTSSessionBegin(const char* params, int* errorCode);
func QTTSSessionBegin(params string) (string, error) {
	var errorCode C.int

	ret := C.QTTSSessionBegin(C.CString(params), &errorCode)
	if errorCode != C.MSP_SUCCESS {
		return "", fmt.Errorf("QTTSSessionBegin failed, error %d", int(errorCode))
	}

	return C.GoString(ret), nil
}

// int MSPAPI QTTSSessionEnd(const char* sessionID, const char* hints);
func QTTSSessionEnd(sessionID string, hints string) error {

	ret := C.QTTSSessionEnd(C.CString(sessionID), C.CString(hints))
	if ret != C.MSP_SUCCESS {
		return fmt.Errorf("QTTSSessionEnd failed, error %d", int(ret))
	}

	return nil
}

// int MSPAPI QTTSSetParam(const char *sessionID, const char *paramName, const char *paramValue);
func QTTSSetParam(sessionID string, paramName string, paramValue string) error {

	ret := C.QTTSSetParam(C.CString(sessionID), C.CString(paramName), C.CString(paramValue))
	if ret != C.MSP_SUCCESS {
		return fmt.Errorf("QTTSSetParam failed, error %d", int(ret))
	}

	return nil
}

// int MSPAPI QTTSTextPut(const char* sessionID, const char* textString, unsigned int textLen, const char* params);
func QTTSTextPut(sessionID string, textString string, params string) error {
	var textLen = C.strlen(C.CString(textString))
	ret := C.QTTSTextPut(C.CString(sessionID), C.CString(textString), C.uint(textLen), C.CString(params))
	if ret != C.MSP_SUCCESS {
		return fmt.Errorf("QTTSTextPut failed, error %d", int(ret))
	}

	return nil
}
