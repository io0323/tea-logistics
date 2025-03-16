import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { ChakraProvider } from '@chakra-ui/react';
import BatchConfigModal from '../BatchConfigModal';
import { BatchType } from '@/types/batch';

describe('BatchConfigModal', () => {
  const mockOnSubmit = jest.fn();
  const mockOnClose = jest.fn();

  const defaultProps = {
    isOpen: true,
    onClose: mockOnClose,
    onSubmit: mockOnSubmit,
    isLoading: false,
  };

  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('正しくレンダリングされること', () => {
    render(
      <ChakraProvider>
        <BatchConfigModal {...defaultProps} />
      </ChakraProvider>
    );

    expect(screen.getByText('バッチ処理の設定')).toBeInTheDocument();
    expect(screen.getByLabelText('処理タイプ')).toBeInTheDocument();
    expect(screen.getByLabelText('スケジュール')).toBeInTheDocument();
    expect(screen.getByLabelText('リトライ回数')).toBeInTheDocument();
    expect(screen.getByLabelText('タイムアウト（秒）')).toBeInTheDocument();
  });

  it('フォームの初期値が正しく設定されていること', () => {
    render(
      <ChakraProvider>
        <BatchConfigModal {...defaultProps} />
      </ChakraProvider>
    );

    expect(screen.getByLabelText('処理タイプ')).toHaveValue(BatchType.STOCK_CHECK);
    expect(screen.getByLabelText('スケジュール')).toHaveValue('');
    expect(screen.getByLabelText('リトライ回数')).toHaveValue('3');
    expect(screen.getByLabelText('タイムアウト（秒）')).toHaveValue('300');
  });

  it('無効なcron形式でエラーが表示されること', async () => {
    render(
      <ChakraProvider>
        <BatchConfigModal {...defaultProps} />
      </ChakraProvider>
    );

    const scheduleInput = screen.getByLabelText('スケジュール');
    fireEvent.change(scheduleInput, { target: { value: 'invalid cron' } });
    fireEvent.click(screen.getByText('設定'));

    await waitFor(() => {
      expect(screen.getByText('無効なcron形式です')).toBeInTheDocument();
    });
  });

  it('タイムアウトがリトライ回数に対して小さすぎる場合にエラーが表示されること', async () => {
    render(
      <ChakraProvider>
        <BatchConfigModal {...defaultProps} />
      </ChakraProvider>
    );

    const retryInput = screen.getByLabelText('リトライ回数');
    const timeoutInput = screen.getByLabelText('タイムアウト（秒）');

    fireEvent.change(retryInput, { target: { value: '5' } });
    fireEvent.change(timeoutInput, { target: { value: '200' } });
    fireEvent.click(screen.getByText('設定'));

    await waitFor(() => {
      expect(screen.getByText('タイムアウト時間はリトライ回数 × 60秒以上に設定してください')).toBeInTheDocument();
    });
  });

  it('正常な入力で送信が成功すること', async () => {
    render(
      <ChakraProvider>
        <BatchConfigModal {...defaultProps} />
      </ChakraProvider>
    );

    const typeSelect = screen.getByLabelText('処理タイプ');
    const scheduleInput = screen.getByLabelText('スケジュール');
    const retryInput = screen.getByLabelText('リトライ回数');
    const timeoutInput = screen.getByLabelText('タイムアウト（秒）');

    fireEvent.change(typeSelect, { target: { value: BatchType.STOCK_CHECK } });
    fireEvent.change(scheduleInput, { target: { value: '*/30 * * * *' } });
    fireEvent.change(retryInput, { target: { value: '3' } });
    fireEvent.change(timeoutInput, { target: { value: '300' } });
    fireEvent.click(screen.getByText('設定'));

    await waitFor(() => {
      expect(mockOnSubmit).toHaveBeenCalledWith({
        type: BatchType.STOCK_CHECK,
        schedule: '*/30 * * * *',
        retryCount: 3,
        timeout: 300,
      });
    });
  });

  it('キャンセルボタンでフォームがリセットされること', () => {
    render(
      <ChakraProvider>
        <BatchConfigModal {...defaultProps} />
      </ChakraProvider>
    );

    const scheduleInput = screen.getByLabelText('スケジュール');
    fireEvent.change(scheduleInput, { target: { value: '*/30 * * * *' } });
    fireEvent.click(screen.getByText('キャンセル'));

    expect(mockOnClose).toHaveBeenCalled();
    expect(scheduleInput).toHaveValue('');
  });
}); 