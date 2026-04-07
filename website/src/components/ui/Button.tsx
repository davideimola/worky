import { type ButtonHTMLAttributes, forwardRef } from 'react';
import { cva, type VariantProps } from 'class-variance-authority';
import { cn } from '@/lib/utils';

const buttonVariants = cva(
  // base
  'inline-flex items-center justify-center gap-2 rounded-md font-mono text-sm font-medium transition-all duration-150 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-teal focus-visible:ring-offset-2 focus-visible:ring-offset-bg disabled:pointer-events-none disabled:opacity-50',
  {
    variants: {
      variant: {
        primary:
          'bg-teal text-midnight hover:bg-teal-deep hover:text-white hover:-translate-y-px active:translate-y-0',
        ghost:
          'text-slate border border-border-mid bg-transparent hover:border-white/40 hover:text-white hover:-translate-y-px active:translate-y-0',
        outline:
          'border border-teal-border text-teal bg-teal-ghost hover:bg-teal/20 hover:-translate-y-px',
        link:
          'text-teal underline-offset-4 hover:underline p-0 h-auto',
      },
      size: {
        sm:      'h-8 px-3 text-xs',
        default: 'h-10 px-5 py-3',
        lg:      'h-12 px-7 text-base',
        icon:    'h-9 w-9 p-0',
      },
    },
    defaultVariants: {
      variant: 'primary',
      size: 'default',
    },
  }
);

export interface ButtonProps
  extends ButtonHTMLAttributes<HTMLButtonElement>,
    VariantProps<typeof buttonVariants> {
  asChild?: boolean;
}

const Button = forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant, size, ...props }, ref) => {
    return (
      <button
        ref={ref}
        className={cn(buttonVariants({ variant, size }), className)}
        {...props}
      />
    );
  }
);

Button.displayName = 'Button';

export { Button, buttonVariants };
