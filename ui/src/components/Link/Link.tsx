import NextLink from 'next/link';

import { IconArrowTopRightOnSquare, IconChevron } from '@/icons';
import classNames from '@/utils/classnames';

interface LinkProps {
  internalNavigation?: boolean;
  children: React.ReactNode;
  className?: string;
  href: string;
}

const defaultLinkStyles =
  'text-indigo-400 hover:decoration-indigo-400 decoration-transparent decoration-2 underline underline-offset-4 cursor-pointer transition-color duration-300 flex items-center gap-1';

export default function Link({
  href,
  children,
  className,
  internalNavigation = false,
}: LinkProps) {
  if (internalNavigation) {
    return (
      <NextLink href={href} className={classNames(className, defaultLinkStyles)}>
        {children}
        <IconChevron className="-rotate-90" />
      </NextLink>
    );
  }
  return (
    <a
      className={classNames(className, defaultLinkStyles)}
      target="_blank"
      href={href}
    >
      {children}
      {<IconArrowTopRightOnSquare />}
    </a>
  );
}
