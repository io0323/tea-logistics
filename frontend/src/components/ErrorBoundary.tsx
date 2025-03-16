import { Component, ErrorInfo, ReactNode } from 'react';
import {
  Box,
  Heading,
  Text,
  Button,
  VStack,
  Code,
  useToast,
} from '@chakra-ui/react';

interface Props {
  children: ReactNode;
  onError?: (error: Error, errorInfo: ErrorInfo) => void;
}

interface State {
  hasError: boolean;
  error: Error | null;
  errorInfo: ErrorInfo | null;
}

/**
 * エラーバウンダリコンポーネント
 * 子コンポーネントでエラーが発生した場合にフォールバックUIを表示します
 */
export class ErrorBoundary extends Component<Props, State> {
  public state: State = {
    hasError: false,
    error: null,
    errorInfo: null,
  };

  public static getDerivedStateFromError(error: Error): State {
    return {
      hasError: true,
      error,
      errorInfo: null,
    };
  }

  public componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    this.logError(error, errorInfo);
    if (this.props.onError) {
      this.props.onError(error, errorInfo);
    }
  }

  private logError = (error: Error, errorInfo: ErrorInfo) => {
    // 開発環境でのみコンソールにエラーを出力
    if (process.env.NODE_ENV === 'development') {
      console.error('エラーが発生しました:', {
        error: {
          name: error.name,
          message: error.message,
          stack: error.stack,
        },
        componentStack: errorInfo.componentStack,
        timestamp: new Date().toISOString(),
      });
    }

    // 本番環境では外部のエラー追跡サービスにエラーを送信
    if (process.env.NODE_ENV === 'production') {
      // エラー追跡サービスへの送信処理
      // 例: Sentry, LogRocket, etc.
      this.sendErrorToTrackingService(error, errorInfo);
    }
  };

  private sendErrorToTrackingService = (error: Error, errorInfo: ErrorInfo) => {
    // TODO: エラー追跡サービスの実装
    // 例: Sentryを使用する場合
    // Sentry.captureException(error, {
    //   extra: {
    //     componentStack: errorInfo.componentStack,
    //   },
    // });
  };

  private handleRetry = () => {
    this.setState({
      hasError: false,
      error: null,
      errorInfo: null,
    });
  };

  public render() {
    if (this.state.hasError) {
      return (
        <Box p={6}>
          <VStack spacing={4} align="start">
            <Heading size="lg" color="red.500">
              エラーが発生しました
            </Heading>
            <Text>
              申し訳ありません。予期せぬエラーが発生しました。
            </Text>
            {process.env.NODE_ENV === 'development' && this.state.error && (
              <Code p={4} borderRadius="md" bg="gray.50" width="100%">
                {this.state.error.toString()}
                {this.state.errorInfo?.componentStack}
              </Code>
            )}
            <Button
              colorScheme="blue"
              onClick={this.handleRetry}
            >
              再試行
            </Button>
          </VStack>
        </Box>
      );
    }

    return this.props.children;
  }
} 