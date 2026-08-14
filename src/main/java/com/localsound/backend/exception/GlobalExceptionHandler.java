package com.localsound.backend.exception;

import java.util.Map;

import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.MethodArgumentNotValidException;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.RestControllerAdvice;

@RestControllerAdvice
public class GlobalExceptionHandler {

	@ExceptionHandler(MethodArgumentNotValidException.class)
	public ResponseEntity<Map<String, String>> handleValidation(
		MethodArgumentNotValidException exception
	) {
		String message = exception.getBindingResult()
			.getFieldError()
			.getDefaultMessage();
		
		return ResponseEntity
			.badRequest()
			.body(Map.of("message", message));
	}
}