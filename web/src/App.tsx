import { useState } from 'react';
import { ThemeProvider } from './contexts';
import ComponentShowcase from './pages/ComponentShowcase';
import TableDemo from './pages/TableDemo';
import { Button } from './components/ui';

function App() {
  const [showTableDemo, setShowTableDemo] = useState(false);

  return (
    <ThemeProvider defaultTheme="modern-blue">
      <div style={{ position: 'fixed', top: '1rem', right: '1rem', zIndex: 1000 }}>
        <Button
          variant="secondary"
          size="small"
          onClick={() => setShowTableDemo(!showTableDemo)}
        >
          {showTableDemo ? '返回组件展示' : '查看表格演示'}
        </Button>
      </div>
      {showTableDemo ? <TableDemo /> : <ComponentShowcase />}
    </ThemeProvider>
  );
}

export default App;
