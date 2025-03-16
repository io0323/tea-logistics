import {
  Modal,
  ModalOverlay,
  ModalContent,
  ModalHeader,
  ModalBody,
  ModalFooter,
  Button,
  FormControl,
  FormLabel,
  Select,
  Input,
  VStack,
  useToast,
  Text,
  NumberInput,
  NumberInputField,
  NumberInputStepper,
  NumberIncrementStepper,
  NumberDecrementStepper,
  FormHelperText,
  FormErrorMessage,
} from '@chakra-ui/react';
import { useState, useCallback, useMemo } from 'react';
import { BatchType, BatchConfig } from '@/types/batch';
import * as cronParser from 'cron-parser';

/**
 * cronスケジュールの検証
 */
const validateCronSchedule = (schedule: string): boolean => {
  if (!schedule) return true;
  try {
    cronParser.parseExpression(schedule);
    return true;
  } catch (error) {
    return false;
  }
};

interface BatchConfigModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (config: BatchConfig) => void;
  isLoading?: boolean;
}

/**
 * バッチ処理の設定モーダル
 */
export default function BatchConfigModal({
  isOpen,
  onClose,
  onSubmit,
  isLoading,
}: BatchConfigModalProps) {
  const toast = useToast();
  const [type, setType] = useState<BatchType>(BatchType.STOCK_CHECK);
  const [schedule, setSchedule] = useState('');
  const [retryCount, setRetryCount] = useState(3);
  const [timeout, setTimeout] = useState(300);
  const [errors, setErrors] = useState<Record<string, string>>({});

  const getTypeLabel = useMemo(() => {
    const labels: Record<BatchType, string> = {
      [BatchType.STOCK_CHECK]: '在庫確認',
      [BatchType.DELIVERY_STATUS_UPDATE]: '配送ステータス更新',
      [BatchType.DATA_CLEANUP]: 'データクリーンアップ',
      [BatchType.REPORT_GENERATION]: 'レポート生成',
    };
    return (type: BatchType): string => labels[type];
  }, []);

  const getTypeDescription = useMemo(() => {
    const descriptions: Record<BatchType, string> = {
      [BatchType.STOCK_CHECK]: '在庫数の確認と不整合の検出を行います',
      [BatchType.DELIVERY_STATUS_UPDATE]: '配送ステータスの一括更新を行います',
      [BatchType.DATA_CLEANUP]: '不要なデータの削除を行います',
      [BatchType.REPORT_GENERATION]: '各種レポートの生成を行います',
    };
    return (type: BatchType): string => descriptions[type];
  }, []);

  const resetForm = useCallback(() => {
    setType(BatchType.STOCK_CHECK);
    setSchedule('');
    setRetryCount(3);
    setTimeout(300);
    setErrors({});
  }, []);

  const validateForm = useCallback((): boolean => {
    const newErrors: Record<string, string> = {};

    if (!type) {
      newErrors.type = '処理タイプを選択してください';
    }

    if (schedule) {
      if (!validateCronSchedule(schedule)) {
        newErrors.schedule = '無効なcron形式です';
      } else {
        try {
          const interval = cronParser.parseExpression(schedule);
          const next = interval.next();
          const now = new Date();
          if (next.getTime() < now.getTime()) {
            newErrors.schedule = '過去の時刻は指定できません';
          }
        } catch (error) {
          newErrors.schedule = 'スケジュールの解析に失敗しました';
        }
      }
    }

    if (retryCount < 0 || retryCount > 10) {
      newErrors.retryCount = 'リトライ回数は0から10の間で指定してください';
    }

    if (timeout < 60 || timeout > 3600) {
      newErrors.timeout = 'タイムアウトは60秒から3600秒の間で指定してください';
    }

    if (timeout < retryCount * 60) {
      newErrors.timeout = 'タイムアウト時間はリトライ回数 × 60秒以上に設定してください';
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  }, [type, schedule, retryCount, timeout]);

  const handleSubmit = useCallback(() => {
    if (!validateForm()) {
      toast({
        title: '入力内容に誤りがあります',
        status: 'error',
        duration: 3000,
        isClosable: true,
      });
      return;
    }

    try {
      const config: BatchConfig = {
        type,
        schedule: schedule || undefined,
        retryCount,
        timeout,
      };

      onSubmit(config);
    } catch (error) {
      toast({
        title: '設定の保存に失敗しました',
        description: 'もう一度お試しください',
        status: 'error',
        duration: 5000,
        isClosable: true,
      });
    }
  }, [type, schedule, retryCount, timeout, validateForm, onSubmit, toast]);

  const handleClose = useCallback(() => {
    resetForm();
    onClose();
  }, [resetForm, onClose]);

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      size="xl"
      aria-labelledby="batch-config-modal-title"
      aria-describedby="batch-config-modal-desc"
      returnFocusOnClose={true}
      closeOnOverlayClick={!isLoading}
    >
      <ModalOverlay />
      <ModalContent>
        <ModalHeader id="batch-config-modal-title">バッチ処理の設定</ModalHeader>
        <ModalBody>
          <Text id="batch-config-modal-desc" mb={4} fontSize="sm">
            バッチ処理の詳細設定を行います。必要な項目を入力してください。
          </Text>
          <VStack spacing={4}>
            <FormControl
              isInvalid={!!errors.type}
              isRequired
              aria-describedby="type-description"
            >
              <FormLabel htmlFor="batch-type">処理タイプ</FormLabel>
              <Select
                id="batch-type"
                value={type}
                onChange={(e) => setType(e.target.value as BatchType)}
                aria-invalid={!!errors.type}
              >
                {Object.values(BatchType).map((t) => (
                  <option key={t} value={t}>
                    {getTypeLabel(t)}
                  </option>
                ))}
              </Select>
              <FormHelperText id="type-description">
                {getTypeDescription(type)}
              </FormHelperText>
              <FormErrorMessage>{errors.type}</FormErrorMessage>
            </FormControl>

            <FormControl
              isInvalid={!!errors.schedule}
              aria-describedby="schedule-description"
            >
              <FormLabel htmlFor="batch-schedule">スケジュール</FormLabel>
              <Input
                id="batch-schedule"
                value={schedule}
                onChange={(e) => setSchedule(e.target.value)}
                placeholder="*/30 * * * *"
                aria-invalid={!!errors.schedule}
                onKeyPress={(e) => {
                  if (e.key === 'Enter') {
                    handleSubmit();
                  }
                }}
              />
              <FormHelperText id="schedule-description">
                cron形式で指定してください（空欄の場合は即時実行）
              </FormHelperText>
              <FormErrorMessage>{errors.schedule}</FormErrorMessage>
            </FormControl>

            <FormControl
              isInvalid={!!errors.retryCount}
              isRequired
              aria-describedby="retry-description"
            >
              <FormLabel htmlFor="batch-retry">リトライ回数</FormLabel>
              <NumberInput
                id="batch-retry"
                value={retryCount}
                onChange={(_, value) => setRetryCount(value)}
                min={0}
                max={10}
                aria-invalid={!!errors.retryCount}
              >
                <NumberInputField />
                <NumberInputStepper>
                  <NumberIncrementStepper />
                  <NumberDecrementStepper />
                </NumberInputStepper>
              </NumberInput>
              <FormHelperText id="retry-description">
                エラー発生時のリトライ回数を指定してください
              </FormHelperText>
              <FormErrorMessage>{errors.retryCount}</FormErrorMessage>
            </FormControl>

            <FormControl
              isInvalid={!!errors.timeout}
              isRequired
              aria-describedby="timeout-description"
            >
              <FormLabel htmlFor="batch-timeout">タイムアウト（秒）</FormLabel>
              <NumberInput
                id="batch-timeout"
                value={timeout}
                onChange={(_, value) => setTimeout(value)}
                min={60}
                max={3600}
                aria-invalid={!!errors.timeout}
              >
                <NumberInputField />
                <NumberInputStepper>
                  <NumberIncrementStepper />
                  <NumberDecrementStepper />
                </NumberInputStepper>
              </NumberInput>
              <FormHelperText id="timeout-description">
                処理のタイムアウト時間を指定してください
              </FormHelperText>
              <FormErrorMessage>{errors.timeout}</FormErrorMessage>
            </FormControl>
          </VStack>
        </ModalBody>

        <ModalFooter>
          <Button
            variant="ghost"
            mr={3}
            onClick={handleClose}
            isDisabled={isLoading}
            aria-label="キャンセル"
          >
            キャンセル
          </Button>
          <Button
            colorScheme="blue"
            onClick={handleSubmit}
            isLoading={isLoading}
            aria-label="設定を保存"
          >
            設定
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  );
} 