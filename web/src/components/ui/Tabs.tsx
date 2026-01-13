import { Tab } from '@headlessui/react';
import type { ReactNode } from 'react';
import './Tabs.css';

export interface TabItem {
  label: string;
  content: ReactNode;
  disabled?: boolean;
}

export interface TabsProps {
  tabs: TabItem[];
  defaultIndex?: number;
  onChange?: (index: number) => void;
  className?: string;
}

export function Tabs({ tabs, defaultIndex = 0, onChange, className = '' }: TabsProps) {
  return (
    <Tab.Group defaultIndex={defaultIndex} onChange={onChange}>
      <div className={`ui-tabs ${className}`}>
        <Tab.List className="ui-tabs-list">
          {tabs.map((tab, index) => (
            <Tab
              key={index}
              disabled={tab.disabled}
              className="ui-tab"
            >
              {({ selected }) => (
                <span className={`ui-tab-label ${selected ? 'ui-tab-label--selected' : ''}`}>
                  {tab.label}
                </span>
              )}
            </Tab>
          ))}
        </Tab.List>

        <Tab.Panels className="ui-tabs-panels">
          {tabs.map((tab, index) => (
            <Tab.Panel key={index} className="ui-tabs-panel">
              {tab.content}
            </Tab.Panel>
          ))}
        </Tab.Panels>
      </div>
    </Tab.Group>
  );
}
