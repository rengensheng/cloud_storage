# Modern UI Component Library

A comprehensive, theme-switchable UI component library built with React, TypeScript, Framer Motion, HeadlessUI, and Lucide React. Designed following modern UI design principles with a focus on accessibility, performance, and user experience.

## Features

- **4 Beautiful Themes**: Modern Blue, Warm Sunset, Neo Mint, and Slate Dark
- **Fully Themeable**: CSS variable-based design system for easy customization
- **Accessible**: WCAG AA compliant with proper ARIA labels and keyboard navigation
- **Smooth Animations**: Powered by Framer Motion for fluid transitions
- **Type-Safe**: Built with TypeScript for better developer experience
- **Custom Class Support**: All components accept custom className props
- **No Tailwind Required**: Pure CSS with custom class names (you can still use Tailwind if you want)

## Design System

### Color System
- **Primary**: Main brand color for buttons and key actions
- **Secondary**: Accent color for secondary elements
- **Success**: Green tones for positive feedback
- **Warning**: Orange/yellow tones for caution
- **Danger**: Red tones for errors and destructive actions
- **Neutrals**: Grayscale for backgrounds, borders, and text

### Typography
Based on an 8pt grid system with the following hierarchy:
- **H1**: 40px, Bold - Page titles
- **H2**: 28px, Semi-bold - Section headers
- **H3**: 22px, Semi-bold - Subsection headers
- **Body L**: 18px, Regular - Large body text
- **Body M**: 16px, Regular - Standard body text
- **Caption**: 12px, Regular - Small labels and notes

### Spacing
All spacing follows an 8pt grid system:
- 4px (small controls minimum)
- 8px, 16px, 24px, 32px, 40px, 48px (standard increments)

## Available Components

### Form Components

#### Button
```tsx
import { Button } from './components/ui';

<Button variant="primary" size="medium" onClick={handleClick}>
  Click Me
</Button>

// With icons
<Button leftIcon={<Mail />} variant="secondary">
  Send Email
</Button>

// Loading state
<Button loading>Processing...</Button>
```

**Props:**
- `variant`: 'primary' | 'secondary' | 'tertiary' | 'danger'
- `size`: 'small' | 'medium' | 'large'
- `fullWidth`: boolean
- `leftIcon`, `rightIcon`: ReactNode
- `loading`: boolean
- `className`: string (custom classes)

#### Input
```tsx
import { Input } from './components/ui';
import { Mail } from 'lucide-react';

<Input
  label="Email"
  type="email"
  placeholder="you@example.com"
  leftIcon={<Mail />}
  status="default"
  helperText="We'll never share your email"
  className="custom-class"
/>
```

**Props:**
- `label`: string
- `status`: 'default' | 'error' | 'success'
- `helperText`: string
- `leftIcon`, `rightIcon`: ReactNode
- `className`: string

#### Textarea
```tsx
import { Textarea } from './components/ui';

<Textarea
  label="Message"
  placeholder="Enter your message"
  resize="vertical"
  helperText="Maximum 500 characters"
  className="custom-class"
/>
```

**Props:**
- `label`: string
- `status`: 'default' | 'error' | 'success'
- `helperText`: string
- `resize`: 'none' | 'vertical' | 'horizontal' | 'both'
- `className`: string

### Selection Components

#### Checkbox
```tsx
import { Checkbox } from './components/ui';

<Checkbox
  checked={checked}
  onChange={setChecked}
  label="I agree to the terms"
  className="custom-class"
/>
```

#### Radio Group
```tsx
import { RadioGroup } from './components/ui';

const options = [
  { value: '1', label: 'Option 1' },
  { value: '2', label: 'Option 2' }
];

<RadioGroup
  label="Choose an option"
  value={value}
  onChange={setValue}
  options={options}
  className="custom-class"
/>
```

#### Switch
```tsx
import { Switch } from './components/ui';

<Switch
  checked={enabled}
  onChange={setEnabled}
  label="Enable notifications"
  className="custom-class"
/>
```

#### Select
```tsx
import { Select } from './components/ui';

const options = [
  { value: 'react', label: 'React' },
  { value: 'vue', label: 'Vue' }
];

<Select
  label="Framework"
  value={value}
  onChange={setValue}
  options={options}
  placeholder="Select a framework"
  className="custom-class"
/>
```

### Layout Components

#### Card
```tsx
import { Card } from './components/ui';

<Card
  padding="medium"
  shadow="medium"
  hover={true}
  className="custom-class"
>
  <h3>Card Title</h3>
  <p>Card content goes here</p>
</Card>
```

**Props:**
- `padding`: 'none' | 'small' | 'medium' | 'large'
- `shadow`: 'none' | 'small' | 'medium' | 'large'
- `hover`: boolean (adds hover effect)
- `className`: string

#### Container
```tsx
import { Container } from './components/ui';

<Container maxWidth="large" padding center>
  <YourContent />
</Container>
```

**Props:**
- `maxWidth`: 'small' | 'medium' | 'large' | 'full'
- `padding`: boolean
- `center`: boolean
- `className`: string

### Complex Components

#### Dialog (Modal)
```tsx
import { Dialog } from './components/ui';

<Dialog
  open={isOpen}
  onClose={() => setIsOpen(false)}
  title="Welcome"
  size="medium"
  className="custom-class"
>
  <p>Dialog content</p>
  <Button onClick={() => setIsOpen(false)}>Close</Button>
</Dialog>
```

**Props:**
- `open`: boolean
- `onClose`: () => void
- `title`: string
- `size`: 'small' | 'medium' | 'large'
- `className`: string

#### Dropdown Menu
```tsx
import { Dropdown } from './components/ui';
import { Settings, User, LogOut } from 'lucide-react';

<Dropdown trigger={<Button>Menu</Button>} className="custom-class">
  <Dropdown.Item icon={<User />} onClick={handleProfile}>
    Profile
  </Dropdown.Item>
  <Dropdown.Item icon={<Settings />} onClick={handleSettings}>
    Settings
  </Dropdown.Item>
  <Dropdown.Item icon={<LogOut />} onClick={handleLogout}>
    Logout
  </Dropdown.Item>
</Dropdown>
```

#### Tabs
```tsx
import { Tabs } from './components/ui';

const tabs = [
  { label: 'Tab 1', content: <div>Content 1</div> },
  { label: 'Tab 2', content: <div>Content 2</div> }
];

<Tabs
  tabs={tabs}
  defaultIndex={0}
  onChange={(index) => console.log(index)}
  className="custom-class"
/>
```

## Theme System

### Using the ThemeProvider

Wrap your app with the `ThemeProvider`:

```tsx
import { ThemeProvider } from './contexts';

function App() {
  return (
    <ThemeProvider defaultTheme="modern-blue">
      <YourApp />
    </ThemeProvider>
  );
}
```

### Switching Themes

```tsx
import { useTheme } from './contexts';

function ThemeSwitcher() {
  const { theme, setTheme } = useTheme();

  return (
    <Button onClick={() => setTheme('warm-sunset')}>
      Switch to Warm Sunset
    </Button>
  );
}
```

### Available Themes
- `modern-blue` - Clean, professional blue theme
- `warm-sunset` - Warm, inviting orange/red theme
- `neo-mint` - Fresh, modern teal/purple theme
- `slate-dark` - Dark mode with blue accents

### Creating Custom Themes

Add a new theme class in `globals.css`:

```css
.theme-custom {
  --color-primary: #your-color;
  --color-primary-hover: #your-hover-color;
  /* ... other variables */
}
```

## Customization

### Using Custom Classes

All components accept a `className` prop that allows you to add your own styles:

```tsx
<Button className="my-custom-button">
  Custom Styled Button
</Button>
```

### CSS Variables

You can override any design token:

```css
:root {
  --color-primary: #your-brand-color;
  --radius-md: 16px;
  --spacing-4: 20px;
}
```

## Accessibility

All components are built with accessibility in mind:

- ✅ WCAG AA color contrast ratios
- ✅ Keyboard navigation support
- ✅ Screen reader compatible
- ✅ Proper ARIA labels
- ✅ Focus indicators
- ✅ Minimum touch target sizes (44x44px)

## Browser Support

- Chrome (latest)
- Firefox (latest)
- Safari (latest)
- Edge (latest)

## Technologies

- **React 19** - UI framework
- **TypeScript** - Type safety
- **Framer Motion** - Animations
- **HeadlessUI** - Accessible component primitives
- **Lucide React** - Icon library
- **Vite** - Build tool

## Development

```bash
# Install dependencies
pnpm install

# Start dev server
pnpm dev

# Build for production
pnpm build
```

## License

MIT

## Credits

Built following modern UI design principles with inspiration from Material Design, Ant Design, and Chakra UI.
