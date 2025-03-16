import { render, screen, fireEvent } from '@testing-library/react';
import { ChakraProvider } from '@chakra-ui/react';
import { ErrorBoundary } from '../ErrorBoundary';

const ThrowError = () => {
  throw new Error('テストエラー');
};

describe('ErrorBoundary', () => {
  beforeEach(() => {
    jest.spyOn(console, 'error').mockImplementation(() => {});
  });

  afterEach(() => {
    jest.restoreAllMocks();
  });

  it('子コンポーネントが正常な場合、そのまま表示されること', () => {
    render(
      <ChakraProvider>
        <ErrorBoundary>
          <div>正常なコンテンツ</div>
        </ErrorBoundary>
      </ChakraProvider>
    );

    expect(screen.getByText('正常なコンテンツ')).toBeInTheDocument();
  });

  it('エラーが発生した場合、フォールバックUIが表示されること', () => {
    render(
      <ChakraProvider>
        <ErrorBoundary>
          <ThrowError />
        </ErrorBoundary>
      </ChakraProvider>
    );

    expect(screen.getByText('エラーが発生しました')).toBeInTheDocument();
    expect(screen.getByText('申し訳ありません。予期せぬエラーが発生しました。')).toBeInTheDocument();
    expect(screen.getByText('テストエラー')).toBeInTheDocument();
  });

  it('再試行ボタンをクリックすると、エラー状態がリセットされること', () => {
    const { rerender } = render(
      <ChakraProvider>
        <ErrorBoundary>
          <ThrowError />
        </ErrorBoundary>
      </ChakraProvider>
    );

    expect(screen.getByText('エラーが発生しました')).toBeInTheDocument();

    fireEvent.click(screen.getByText('再試行'));

    rerender(
      <ChakraProvider>
        <ErrorBoundary>
          <div>正常なコンテンツ</div>
        </ErrorBoundary>
      </ChakraProvider>
    );

    expect(screen.getByText('正常なコンテンツ')).toBeInTheDocument();
  });
}); 