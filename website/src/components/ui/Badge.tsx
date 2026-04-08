import { type HTMLAttributes } from 'react';
import { cva, type VariantProps } from 'class-variance-authority';
import { cn } from '@/lib/utils';

const badgeVariants = cva(
  'inline-flex items-center gap-2 rounded-full border font-mono text-xs transition-colors',
  {
    variants: {
      variant: {
        teal:    'border-teal-border bg-teal-ghost text-teal px-3.5 py-1.5',
        surface: 'border-border bg-surface text-slate px-3 py-1',
        outline: 'border-border-mid text-slate-light px-3 py-1',
      },
    },
    defaultVariants: {
      variant: 'teal',
    },
  }
);

export interface BadgeProps
  extends HTMLAttributes<HTMLSpanElement>,
    VariantProps<typeof badgeVariants> {}

function Badge({ className, variant, ...props }: BadgeProps) {
  return (
    <span className={cn(badgeVariants({ variant }), className)} {...props} />
  );
}

export { Badge, badgeVariants };
