/**
 * 匯入/匯出對話框組件
 */

import { useState, useRef, type ChangeEvent } from 'react';
import { Upload, Download, FileText, AlertCircle, CheckCircle } from 'lucide-react';
import { Button } from '@/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { Alert, AlertDescription } from '@/components/ui/alert';
import {
  parseCSV,
  generateCSVTemplate,
  exportToJSON,
  exportToCSV,
  downloadFile,
  readFile,
  type CSVParseResult,
} from '@/utils/importExport';
import type { Point, CreatePointRequest } from '@/types/datalink';

// =============================================================================
// 匯入對話框
// =============================================================================

interface ImportDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  deviceId: string;
  deviceName?: string;
  onImport: (points: CreatePointRequest[]) => Promise<void>;
}

export function ImportDialog({
  open,
  onOpenChange,
  deviceId,
  deviceName,
  onImport,
}: ImportDialogProps) {
  const [file, setFile] = useState<File | null>(null);
  const [parseResult, setParseResult] = useState<CSVParseResult | null>(null);
  const [importing, setImporting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleFileChange = async (e: ChangeEvent<HTMLInputElement>) => {
    const selectedFile = e.target.files?.[0];
    if (!selectedFile) return;

    setFile(selectedFile);
    setError(null);

    try {
      const content = await readFile(selectedFile);
      const result = parseCSV(content, deviceId);
      setParseResult(result);
    } catch {
      setError('讀取檔案失敗');
      setParseResult(null);
    }
  };

  const handleImport = async () => {
    if (!parseResult?.success || parseResult.points.length === 0) return;

    setImporting(true);
    try {
      await onImport(parseResult.points);
      onOpenChange(false);
      resetState();
    } catch (err) {
      setError(err instanceof Error ? err.message : '匯入失敗');
    } finally {
      setImporting(false);
    }
  };

  const handleDownloadTemplate = () => {
    const template = generateCSVTemplate();
    downloadFile(template, 'points_template.csv', 'text/csv');
  };

  const resetState = () => {
    setFile(null);
    setParseResult(null);
    setError(null);
    if (fileInputRef.current) {
      fileInputRef.current.value = '';
    }
  };

  const handleClose = (open: boolean) => {
    if (!open) {
      resetState();
    }
    onOpenChange(open);
  };

  return (
    <Dialog open={open} onOpenChange={handleClose}>
      <DialogContent className="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <Upload className="h-5 w-5" />
            匯入點位
          </DialogTitle>
          <DialogDescription>
            從 CSV 檔案匯入點位到 {deviceName || '選定設備'}
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-4 py-4">
          {/* 檔案選擇 */}
          <div className="flex items-center gap-2">
            <input
              ref={fileInputRef}
              type="file"
              accept=".csv"
              onChange={handleFileChange}
              className="hidden"
              id="csv-file-input"
            />
            <Button
              variant="outline"
              onClick={() => fileInputRef.current?.click()}
              className="flex-1"
            >
              <FileText className="mr-2 h-4 w-4" />
              {file ? file.name : '選擇 CSV 檔案'}
            </Button>
            <Button variant="ghost" size="sm" onClick={handleDownloadTemplate}>
              下載範本
            </Button>
          </div>

          {/* 解析結果 */}
          {parseResult && (
            <div className="space-y-2">
              {parseResult.success ? (
                <Alert className="border-green-200 bg-green-50">
                  <CheckCircle className="h-4 w-4 text-green-600" />
                  <AlertDescription className="text-green-700">
                    成功解析 {parseResult.points.length} 個點位
                  </AlertDescription>
                </Alert>
              ) : (
                <Alert variant="destructive">
                  <AlertCircle className="h-4 w-4" />
                  <AlertDescription>
                    解析錯誤:
                    <ul className="mt-1 list-disc list-inside">
                      {parseResult.errors.slice(0, 5).map((err, i) => (
                        <li key={i}>
                          第 {err.row} 行{err.column ? ` (${err.column})` : ''}: {err.message}
                        </li>
                      ))}
                      {parseResult.errors.length > 5 && (
                        <li>...還有 {parseResult.errors.length - 5} 個錯誤</li>
                      )}
                    </ul>
                  </AlertDescription>
                </Alert>
              )}

              {parseResult.warnings.length > 0 && (
                <Alert className="border-yellow-200 bg-yellow-50">
                  <AlertCircle className="h-4 w-4 text-yellow-600" />
                  <AlertDescription className="text-yellow-700">
                    警告: {parseResult.warnings.join('; ')}
                  </AlertDescription>
                </Alert>
              )}
            </div>
          )}

          {/* 錯誤訊息 */}
          {error && (
            <Alert variant="destructive">
              <AlertCircle className="h-4 w-4" />
              <AlertDescription>{error}</AlertDescription>
            </Alert>
          )}
        </div>

        <DialogFooter>
          <Button variant="outline" onClick={() => handleClose(false)}>
            取消
          </Button>
          <Button
            onClick={handleImport}
            disabled={!parseResult?.success || importing}
          >
            {importing ? '匯入中...' : `匯入 ${parseResult?.points.length || 0} 個點位`}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}

// =============================================================================
// 匯出對話框
// =============================================================================

interface ExportDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  points: Point[];
  deviceId?: string;
  deviceName?: string;
}

type ExportFormat = 'csv' | 'json';

export function ExportDialog({
  open,
  onOpenChange,
  points,
  deviceId,
  deviceName,
}: ExportDialogProps) {
  const [format, setFormat] = useState<ExportFormat>('csv');

  const handleExport = () => {
    const timestamp = new Date().toISOString().slice(0, 10);
    const baseName = deviceName ? `${deviceName}_points` : 'points';

    if (format === 'csv') {
      const content = exportToCSV(points);
      downloadFile(content, `${baseName}_${timestamp}.csv`, 'text/csv');
    } else {
      const content = exportToJSON(points, { deviceId, deviceName });
      downloadFile(content, `${baseName}_${timestamp}.json`, 'application/json');
    }

    onOpenChange(false);
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[400px]">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <Download className="h-5 w-5" />
            匯出點位
          </DialogTitle>
          <DialogDescription>
            匯出 {points.length} 個點位{deviceName ? ` (${deviceName})` : ''}
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-4 py-4">
          {/* 格式選擇 */}
          <div className="space-y-2">
            <label className="text-sm font-medium">匯出格式</label>
            <div className="flex gap-4">
              <label className="flex items-center gap-2 cursor-pointer">
                <input
                  type="radio"
                  name="format"
                  value="csv"
                  checked={format === 'csv'}
                  onChange={() => setFormat('csv')}
                  className="h-4 w-4"
                />
                <span>CSV (適合 Excel 開啟)</span>
              </label>
              <label className="flex items-center gap-2 cursor-pointer">
                <input
                  type="radio"
                  name="format"
                  value="json"
                  checked={format === 'json'}
                  onChange={() => setFormat('json')}
                  className="h-4 w-4"
                />
                <span>JSON (完整設定)</span>
              </label>
            </div>
          </div>

          {/* 預覽資訊 */}
          <div className="rounded-md border p-3 bg-muted/50">
            <div className="text-sm space-y-1">
              <p><strong>點位數量:</strong> {points.length}</p>
              {deviceName && <p><strong>設備:</strong> {deviceName}</p>}
              <p><strong>格式:</strong> {format === 'csv' ? 'CSV 表格' : 'JSON 設定檔'}</p>
            </div>
          </div>
        </div>

        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)}>
            取消
          </Button>
          <Button onClick={handleExport} disabled={points.length === 0}>
            <Download className="mr-2 h-4 w-4" />
            匯出
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
