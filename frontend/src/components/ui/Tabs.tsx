import { useState } from 'react';
import React from 'react';

export interface Tab {
  id: string;
  label: string;
}

export interface TabsProps {
  tabs: Tab[];
  defaultTab: string;
  children: React.ReactNode;
}

export function Tabs({ tabs, defaultTab, children }: TabsProps) {
  const [active, setActive] = useState(defaultTab);

  // Find the active tab panel from children
  const activePanel = React.Children.toArray(children).find(
    (child) => React.isValidElement(child) && child.props.id === active
  );

  return (
    <div>
      <div className="flex border-b border-gray-200 dark:border-gray-700">
        {tabs.map((t) => (
          <button
            key={t.id}
            onClick={() => setActive(t.id)}
            className={`px-4 py-2 border-b-2 font-medium text-sm transition-colors ${
              active === t.id
                ? 'border-nfl-dark text-nfl-dark dark:border-nfl-light dark:text-nfl-light'
                : 'border-transparent text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200'
            }`}
          >
            {t.label}
          </button>
        ))}
      </div>
      <div className="py-4" data-active-tab={active}>
        {activePanel}
      </div>
    </div>
  );
}

export function TabPanel({ id, children }: { id: string; children: React.ReactNode }) {
  return <div data-tab-panel={id}>{children}</div>;
}