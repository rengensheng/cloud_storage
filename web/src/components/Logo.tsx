interface LogoProps {
  size?: number;
  className?: string;
}

export function Logo({ size = 32, className = '' }: LogoProps) {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 512 512"
      fill="none"
      width={size}
      height={size}
      className={className}
    >
      <defs>
        <linearGradient id="logoGradient" x1="0%" y1="0%" x2="100%" y2="100%">
          <stop offset="0%" stopColor="#0ea5e9" stopOpacity={1} />
          <stop offset="50%" stopColor="#0284c7" stopOpacity={1} />
          <stop offset="100%" stopColor="#0c4a6e" stopOpacity={1} />
        </linearGradient>

        <linearGradient id="folderGradient" x1="0%" y1="0%" x2="100%" y2="100%">
          <stop offset="0%" stopColor="#0ea5e9" stopOpacity={1} />
          <stop offset="100%" stopColor="#0284c7" stopOpacity={1} />
        </linearGradient>
      </defs>

      {/* Background circle */}
      <circle cx="256" cy="256" r="256" fill="url(#logoGradient)" />

      {/* Cloud shape */}
      <path
        d="M140 320C106.86 320 80 293.14 80 260C80 229.07 103.39 203.6 133.32 200.46C140.79 155.37 179.61 120 226.67 120C277.95 120 319.74 156.69 327.67 205.12C359.63 206.53 385 232.9 385 265C385 295.93 361.61 321.4 331.68 324.54C324.21 369.63 285.39 405 238.33 405C187.05 405 145.26 368.31 137.33 319.88C127.07 319.13 117.47 316.17 109.33 310.67C100.67 304.8 94.27 296.27 91.07 286.27C87.87 276.27 88.07 265.47 91.67 256.13C95.27 246.8 102 239.2 110.67 234.67C119.33 230.13 129.33 228.8 140 230.67"
        fill="white"
        fillOpacity="0.95"
      />

      {/* Folder/Document icon */}
      <path
        d="M190 240H175C171.13 240 168 243.13 168 247V293C168 296.87 171.13 300 175 300H305C308.87 300 312 296.87 312 293V247C312 243.13 308.87 240 305 240H270V254H190V240Z"
        fill="url(#folderGradient)"
      />

      {/* Lock/Security icon */}
      <path
        d="M240 265V255C240 249.48 244.48 245 250 245H254V265H240Z"
        fill="#0ea5e9"
      />

      {/* Small accent clouds */}
      <ellipse cx="165" cy="175" rx="25" ry="18" fill="white" fillOpacity="0.3" />
      <ellipse cx="345" cy="195" rx="20" ry="15" fill="white" fillOpacity="0.3" />
    </svg>
  );
}
