
import React, { useEffect } from 'react';
import { createPortal } from 'react-dom';

// Simple close icon
const CloseIcon = () => (
  <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
  </svg>
);

export interface SlidePanelProps {
  isOpen: boolean;
  title: string;
  onClose: () => void;
  children: React.ReactNode;
  width?: 'sm' | 'md' | 'lg';
}

export function SlidePanel({
  isOpen,
  title,
  onClose,
  children,
  width = 'md',
}: SlidePanelProps) {
  // Handle ESC key
  useEffect(() => {
    const handleEsc = (e: KeyboardEvent) => {
      if (e.key === 'Escape' && isOpen) {
        onClose();
      }
    };
    window.addEventListener('keydown', handleEsc);
    return () => window.removeEventListener('keydown', handleEsc);
  }, [isOpen, onClose]);

  /* 
  // 移除此段落以避免背景滾動條消失導致的頁面抖動 (跳動)
  useEffect(() => {
    if (isOpen) {
      document.body.style.overflow = 'hidden';
    } else {
      document.body.style.overflow = '';
    }
    return () => {
      document.body.style.overflow = '';
    };
  }, [isOpen]);
  */

  if (!isOpen) return null;

  const widthClass = {
    sm: 'w-80',
    md: 'w-96',
    lg: 'w-[32rem]',
  }[width];

  return createPortal(
    <div className="fixed inset-0 z-50 flex justify-end" data-testid="slide-panel">
      {/* Backdrop */}
      <div 
        className="fixed inset-0 bg-black/20 backdrop-blur-[2px] transition-opacity duration-300"
        onClick={onClose}
      />
      
      {/* Panel */}
      <div 
        className={`
          relative h-full bg-white dark:bg-slate-900 border-l border-slate-200 dark:border-slate-700 shadow-2xl 
          transform transition-transform duration-300 ease-out flex flex-col
          ${widthClass}
        `}
      >
        {/* Header */}
        <div className="h-14 flex items-center justify-between px-6 border-b border-slate-200 dark:border-slate-700 flex-shrink-0">
          <h2 className="text-lg font-semibold text-slate-800 dark:text-slate-100">
            {title}
          </h2>
          <button
            onClick={onClose}
            className="p-2 -mr-2 text-slate-500 hover:text-slate-700 dark:text-slate-400 dark:hover:text-slate-200 rounded-full hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors"
          >
            <CloseIcon />
          </button>
        </div>

        {/* Content */}
        <div className="flex-1 overflow-y-auto p-6">
          {children}
        </div>
      </div>
    </div>,
    document.body
  );
}
