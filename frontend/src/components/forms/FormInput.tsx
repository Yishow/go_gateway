import type { InputHTMLAttributes, ReactNode } from 'react';
import { designSystem } from '../../styles/designSystem';
import FormLabel from './FormLabel';
import FormError from './FormError';
import FormHint from './FormHint';
import FormSuccess from './FormSuccess';

export interface FormInputProps extends Omit<InputHTMLAttributes<HTMLInputElement>, 'className'> {
  /** Label 文字 */
  label?: ReactNode;
  /** 是否為必填欄位 */
  required?: boolean;
  /** 錯誤訊息 */
  error?: ReactNode;
  /** 提示訊息 */
  hint?: ReactNode;
  /** 成功訊息 */
  success?: ReactNode;
  /** 額外的 input className */
  inputClassName?: string;
  /** 額外的 wrapper className */
  className?: string;
  /** 是否顯示錯誤圖示 */
  showErrorIcon?: boolean;
  /** 是否顯示成功圖示 */
  showSuccessIcon?: boolean;
}

/**
 * 統一的表單 Input 元件
 * 整合 Label、Input、Error、Hint、Success 訊息
 * 使用 designSystem.forms.* 樣式
 */
export default function FormInput({
  label,
  required = false,
  error,
  hint,
  success,
  inputClassName = '',
  className = '',
  showErrorIcon = true,
  showSuccessIcon = true,
  id,
  ...inputProps
}: FormInputProps) {
  const inputId = id || `form-input-${Math.random().toString(36).substr(2, 9)}`;
  const hasError = !!error;
  const hasSuccess = !!success && !hasError;

  return (
    <div className={className}>
      {label && (
        <FormLabel htmlFor={inputId} required={required}>
          {label}
        </FormLabel>
      )}
      <input
        {...inputProps}
        id={inputId}
        className={`${designSystem.forms.input.base} ${hasError ? designSystem.forms.input.error : ''} ${inputClassName}`}
        aria-invalid={hasError}
        aria-describedby={
          error || hint || success
            ? `${inputId}-${error ? 'error' : success ? 'success' : 'hint'}`
            : undefined
        }
      />
      {error && (
        <FormError id={`${inputId}-error`} showIcon={showErrorIcon}>
          {error}
        </FormError>
      )}
      {success && !hasError && (
        <FormSuccess id={`${inputId}-success`} showIcon={showSuccessIcon}>
          {success}
        </FormSuccess>
      )}
      {hint && !hasError && !hasSuccess && (
        <FormHint id={`${inputId}-hint`}>
          {hint}
        </FormHint>
      )}
    </div>
  );
}
