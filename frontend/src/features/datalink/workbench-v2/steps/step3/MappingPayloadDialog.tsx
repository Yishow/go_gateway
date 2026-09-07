import * as React from 'react';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '../../../../../components/ui/dialog';
import type { Point, Mapping } from '../../state/types';
import { PayloadPreview } from './PayloadPreview';

export interface MappingPayloadDialogProps {
  point: Point;
  mapping: Mapping;
}

export const MappingPayloadDialog: React.FC<MappingPayloadDialogProps> = ({
  point,
  mapping,
}) => {
  const [open, setOpen] = React.useState(false);
  const stopPropagation = (event: React.MouseEvent) => {
    event.stopPropagation();
  };

  return (
    <td className="p-2 text-center">
      <Dialog open={open} onOpenChange={setOpen}>
        <DialogTrigger asChild>
          <button
            type="button"
            onMouseDown={stopPropagation}
            onClick={(event) => {
              event.stopPropagation();
              setOpen(true);
            }}
            className="rounded border border-slate-800 bg-slate-950/40 px-2 py-1 text-[11px] font-semibold text-slate-300 transition-colors hover:border-slate-700 hover:text-slate-100"
            data-testid={`payload-modal-trigger-${point.id}`}
          >
            JSON
          </button>
        </DialogTrigger>
        <DialogContent
          className="max-w-3xl border-slate-800 bg-slate-950 text-slate-100"
          data-testid={`payload-modal-content-${point.id}`}
        >
          <DialogHeader>
            <DialogTitle className="text-slate-100">{mapping.tag_key || point.name}</DialogTitle>
            <DialogDescription className="text-slate-400">
              {point.name} @ {point.address}
            </DialogDescription>
          </DialogHeader>
          <PayloadPreview point={point} mapping={mapping} />
        </DialogContent>
      </Dialog>
    </td>
  );
};
