import { type FC } from 'react';

interface LoadingProps {
  message?: string;
  size?: 'sm' | 'md' | 'lg';
}

const LoadingSpinner: FC<LoadingProps> = ({ 
  message = 'Loading...', 
  size = 'md' 
}) => {
  const sizeClasses = {
    sm: 'w-6 h-6',
    md: 'w-10 h-10',
    lg: 'w-16 h-16'
  };

  const borderSize = {
    sm: 'border-2',
    md: 'border-4',
    lg: 'border-6'
  };

  return (
    <div className="flex flex-col items-center justify-center min-h-[200px]">
      <div className={`${sizeClasses[size]} ${borderSize[size]} border-t-blue-500 border-r-transparent border-b-transparent border-l-transparent rounded-full animate-spin mb-4`}></div>
      <p className="text-gray-600">{message}</p>
    </div>
  );
};

export default LoadingSpinner;