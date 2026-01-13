import { useState } from 'react';
import { useTheme } from '../contexts';
import type { Theme } from '../types';
import {
  Button,
  Input,
  Textarea,
  Checkbox,
  RadioGroup,
  Switch,
  Select,
  Card,
  Container,
  Dialog,
  Dropdown,
  Tabs
} from '../components/ui';
import {
  Palette,
  Mail,
  Lock,
  User,
  Settings,
  LogOut,
  Trash2,
  Edit
} from 'lucide-react';
import '../styles/globals.css';

const themes: { value: Theme; label: string; description: string }[] = [
  { value: 'modern-blue', label: 'Modern Blue', description: '专业蓝' },
  { value: 'warm-sunset', label: 'Warm Sunset', description: '温暖橙' },
  { value: 'neo-mint', label: 'Neo Mint', description: '清新薄荷' },
  { value: 'slate-dark', label: 'Slate Dark', description: '深色模式' },
  { value: 'purple-dream', label: 'Purple Dream', description: '梦幻紫' },
  { value: 'ocean-breeze', label: 'Ocean Breeze', description: '海洋蓝' },
  { value: 'forest-green', label: 'Forest Green', description: '森林绿' },
  { value: 'rose-gold', label: 'Rose Gold', description: '玫瑰金' },
  { value: 'midnight-purple', label: 'Midnight Purple', description: '午夜紫' },
  { value: 'sakura-pink', label: 'Sakura Pink', description: '樱花粉' },
  { value: 'cyber-neon', label: 'Cyber Neon', description: '赛博霓虹' }
];

const radioOptions = [
  { value: 'option1', label: 'Option 1' },
  { value: 'option2', label: 'Option 2' },
  { value: 'option3', label: 'Option 3' }
];

const selectOptions = [
  { value: 'react', label: 'React' },
  { value: 'vue', label: 'Vue' },
  { value: 'angular', label: 'Angular' },
  { value: 'svelte', label: 'Svelte' }
];

export default function ComponentShowcase() {
  const { theme, setTheme } = useTheme();
  const [checked, setChecked] = useState(false);
  const [switchOn, setSwitchOn] = useState(false);
  const [radioValue, setRadioValue] = useState('option1');
  const [selectValue, setSelectValue] = useState('react');
  const [dialogOpen, setDialogOpen] = useState(false);
  const [email, setEmail] = useState('');
  const [message, setMessage] = useState('');

  return (
    <div style={{ minHeight: '100vh', paddingTop: '2rem', paddingBottom: '2rem' }}>
      <Container maxWidth="large">
        {/* Header */}
        <div style={{ marginBottom: '3rem', textAlign: 'center' }}>
          <h1 className="text-h1" style={{ marginBottom: '0.5rem' }}>
            UI Component Library
          </h1>
          <p className="text-body-l" style={{ color: 'var(--color-text-muted)' }}>
            A modern, theme-switchable component system built with React, Framer Motion, and HeadlessUI
          </p>
        </div>

        {/* Theme Switcher */}
        <Card padding="medium" shadow="medium" style={{ marginBottom: '3rem' }}>
          <div style={{ marginBottom: '1rem' }}>
            <div style={{ display: 'flex', alignItems: 'center', gap: '1rem', marginBottom: '1rem' }}>
              <Palette style={{ color: 'var(--color-primary)' }} />
              <h3 className="text-h3" style={{ margin: 0 }}>
                主题切换器 ({themes.length} 个主题)
              </h3>
            </div>
            <p className="text-body-s" style={{ color: 'var(--color-text-muted)', marginBottom: 0 }}>
              点击下方按钮切换不同主题，体验多样化的视觉风格
            </p>
          </div>
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(150px, 1fr))', gap: '0.75rem' }}>
            {themes.map((t) => (
              <Button
                key={t.value}
                variant={theme === t.value ? 'primary' : 'secondary'}
                size="small"
                onClick={() => setTheme(t.value)}
                style={{ flexDirection: 'column', height: 'auto', padding: '0.5rem' }}
              >
                <span style={{ fontWeight: 600, fontSize: '13px' }}>{t.label}</span>
                <span style={{ fontSize: '11px', opacity: 0.8 }}>{t.description}</span>
              </Button>
            ))}
          </div>
        </Card>

        {/* Components Showcase */}
        <Tabs
          tabs={[
            {
              label: 'Buttons',
              content: (
                <div style={{ display: 'flex', flexDirection: 'column', gap: '2rem' }}>
                  <Card padding="large">
                    <h3 className="text-h3" style={{ marginBottom: '1rem' }}>Button Variants</h3>
                    <div style={{ display: 'flex', gap: '1rem', flexWrap: 'wrap' }}>
                      <Button variant="primary">Primary</Button>
                      <Button variant="secondary">Secondary</Button>
                      <Button variant="tertiary">Tertiary</Button>
                      <Button variant="danger">Danger</Button>
                    </div>
                  </Card>

                  <Card padding="large">
                    <h3 className="text-h3" style={{ marginBottom: '1rem' }}>Button Sizes</h3>
                    <div style={{ display: 'flex', gap: '1rem', alignItems: 'center', flexWrap: 'wrap' }}>
                      <Button size="small">Small</Button>
                      <Button size="medium">Medium</Button>
                      <Button size="large">Large</Button>
                    </div>
                  </Card>

                  <Card padding="large">
                    <h3 className="text-h3" style={{ marginBottom: '1rem' }}>Buttons with Icons</h3>
                    <div style={{ display: 'flex', gap: '1rem', flexWrap: 'wrap' }}>
                      <Button leftIcon={<Mail />}>Send Email</Button>
                      <Button variant="secondary" rightIcon={<Settings />}>Settings</Button>
                      <Button variant="danger" leftIcon={<Trash2 />}>Delete</Button>
                    </div>
                  </Card>

                  <Card padding="large">
                    <h3 className="text-h3" style={{ marginBottom: '1rem' }}>Button States</h3>
                    <div style={{ display: 'flex', gap: '1rem', flexWrap: 'wrap' }}>
                      <Button loading>Loading</Button>
                      <Button disabled>Disabled</Button>
                      <Button fullWidth>Full Width Button</Button>
                    </div>
                  </Card>
                </div>
              )
            },
            {
              label: 'Form Inputs',
              content: (
                <div style={{ display: 'flex', flexDirection: 'column', gap: '2rem' }}>
                  <Card padding="large">
                    <h3 className="text-h3" style={{ marginBottom: '1rem' }}>Text Inputs</h3>
                    <div style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
                      <Input
                        label="Email"
                        type="email"
                        placeholder="Enter your email"
                        value={email}
                        onChange={(e) => setEmail(e.target.value)}
                      />
                      <Input
                        label="Email with Icon"
                        type="email"
                        placeholder="you@example.com"
                        leftIcon={<Mail />}
                      />
                      <Input
                        label="Password"
                        type="password"
                        placeholder="Enter password"
                        leftIcon={<Lock />}
                      />
                      <Input
                        status="error"
                        placeholder="Error state"
                        helperText="This field has an error"
                      />
                      <Input
                        status="success"
                        placeholder="Success state"
                        helperText="Looks good!"
                      />
                    </div>
                  </Card>

                  <Card padding="large">
                    <h3 className="text-h3" style={{ marginBottom: '1rem' }}>Textarea</h3>
                    <Textarea
                      label="Message"
                      placeholder="Enter your message"
                      value={message}
                      onChange={(e) => setMessage(e.target.value)}
                      helperText="Share your thoughts with us"
                    />
                  </Card>
                </div>
              )
            },
            {
              label: 'Selections',
              content: (
                <div style={{ display: 'flex', flexDirection: 'column', gap: '2rem' }}>
                  <Card padding="large">
                    <h3 className="text-h3" style={{ marginBottom: '1rem' }}>Checkbox</h3>
                    <div style={{ display: 'flex', flexDirection: 'column', gap: '0.5rem' }}>
                      <Checkbox
                        checked={checked}
                        onChange={setChecked}
                        label="I agree to the terms and conditions"
                      />
                      <Checkbox label="Subscribe to newsletter" />
                      <Checkbox label="Disabled option" disabled />
                    </div>
                  </Card>

                  <Card padding="large">
                    <h3 className="text-h3" style={{ marginBottom: '1rem' }}>Radio Group</h3>
                    <RadioGroup
                      label="Choose an option"
                      value={radioValue}
                      onChange={setRadioValue}
                      options={radioOptions}
                    />
                  </Card>

                  <Card padding="large">
                    <h3 className="text-h3" style={{ marginBottom: '1rem' }}>Switch</h3>
                    <div style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
                      <Switch
                        checked={switchOn}
                        onChange={setSwitchOn}
                        label="Enable notifications"
                      />
                      <Switch label="Dark mode" />
                      <Switch label="Disabled switch" disabled />
                    </div>
                  </Card>

                  <Card padding="large">
                    <h3 className="text-h3" style={{ marginBottom: '1rem' }}>Select</h3>
                    <Select
                      label="Choose a framework"
                      value={selectValue}
                      onChange={setSelectValue}
                      options={selectOptions}
                      placeholder="Select a framework"
                    />
                  </Card>
                </div>
              )
            },
            {
              label: 'Complex',
              content: (
                <div style={{ display: 'flex', flexDirection: 'column', gap: '2rem' }}>
                  <Card padding="large">
                    <h3 className="text-h3" style={{ marginBottom: '1rem' }}>Dialog (Modal)</h3>
                    <Button onClick={() => setDialogOpen(true)}>Open Dialog</Button>
                    <Dialog
                      open={dialogOpen}
                      onClose={() => setDialogOpen(false)}
                      title="Welcome to our platform"
                    >
                      <div style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
                        <p className="text-body-m">
                          This is a modal dialog component built with HeadlessUI and Framer Motion.
                          It features smooth animations and is fully accessible.
                        </p>
                        <Input
                          label="Your Name"
                          placeholder="John Doe"
                          leftIcon={<User />}
                        />
                        <div style={{ display: 'flex', gap: '0.5rem', justifyContent: 'flex-end' }}>
                          <Button variant="secondary" onClick={() => setDialogOpen(false)}>
                            Cancel
                          </Button>
                          <Button onClick={() => setDialogOpen(false)}>Submit</Button>
                        </div>
                      </div>
                    </Dialog>
                  </Card>

                  <Card padding="large">
                    <h3 className="text-h3" style={{ marginBottom: '1rem' }}>Dropdown Menu</h3>
                    <Dropdown trigger={<Button>Open Menu</Button>}>
                      <Dropdown.Item icon={<User />} onClick={() => alert('Profile')}>
                        Profile
                      </Dropdown.Item>
                      <Dropdown.Item icon={<Settings />} onClick={() => alert('Settings')}>
                        Settings
                      </Dropdown.Item>
                      <Dropdown.Item icon={<Edit />} onClick={() => alert('Edit')}>
                        Edit
                      </Dropdown.Item>
                      <Dropdown.Item icon={<LogOut />} onClick={() => alert('Logout')}>
                        Logout
                      </Dropdown.Item>
                    </Dropdown>
                  </Card>

                  <Card padding="large">
                    <h3 className="text-h3" style={{ marginBottom: '1rem' }}>Cards</h3>
                    <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(250px, 1fr))', gap: '1rem' }}>
                      <Card padding="medium" shadow="small">
                        <h4 className="text-h3" style={{ marginBottom: '0.5rem' }}>Card 1</h4>
                        <p className="text-body-s">A simple card with small shadow</p>
                      </Card>
                      <Card padding="medium" shadow="medium">
                        <h4 className="text-h3" style={{ marginBottom: '0.5rem' }}>Card 2</h4>
                        <p className="text-body-s">A card with medium shadow</p>
                      </Card>
                      <Card padding="medium" shadow="large" hover>
                        <h4 className="text-h3" style={{ marginBottom: '0.5rem' }}>Hover Card</h4>
                        <p className="text-body-s">Hover over me!</p>
                      </Card>
                    </div>
                  </Card>
                </div>
              )
            }
          ]}
        />

        {/* Footer */}
        <div style={{ marginTop: '3rem', textAlign: 'center', padding: '2rem 0', borderTop: '1px solid var(--color-border)' }}>
          <p className="text-caption">
            Built with React, TypeScript, Framer Motion, HeadlessUI, and Lucide React
          </p>
        </div>
      </Container>
    </div>
  );
}
