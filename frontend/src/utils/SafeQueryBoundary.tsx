import * as React from 'react';
import { QueryClient, QueryClientContext, QueryClientProvider } from '@tanstack/react-query';

const fallbackQueryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: false,
      refetchOnWindowFocus: false,
    },
  },
});

export interface SafeQueryBoundaryProps {
  children: React.ReactNode;
}

export const SafeQueryBoundary: React.FC<SafeQueryBoundaryProps> = ({ children }) => {
  const client = React.useContext(QueryClientContext);
  if (!client) {
    return <QueryClientProvider client={fallbackQueryClient}>{children}</QueryClientProvider>;
  }
  return <>{children}</>;
};
