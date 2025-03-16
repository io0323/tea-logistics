'use client';

import { useState } from 'react';
import {
  Box,
  VStack,
  FormControl,
  FormLabel,
  Input,
  Select,
  Switch,
  Button,
  useToast,
  useColorModeValue,
} from '@chakra-ui/react';
import { useSettings } from '@/hooks/useSettings';

interface SettingsFormProps {
  initialSettings: any;
}

/**
 * 設定フォームコンポーネント
 */
export function SettingsForm({ initialSettings }: SettingsFormProps) {
  const toast = useToast();
  const { updateSettings } = useSettings();
  const bgColor = useColorModeValue('white', 'gray.700');

  const [formData, setFormData] = useState({
    companyName: initialSettings?.companyName || '',
    email: initialSettings?.email || '',
    notificationEnabled: initialSettings?.notificationEnabled || false,
    language: initialSettings?.language || 'ja',
    timezone: initialSettings?.timezone || 'Asia/Tokyo',
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await updateSettings(formData);
      toast({
        title: '設定を更新しました',
        status: 'success',
        duration: 3000,
        isClosable: true,
      });
    } catch (error) {
      toast({
        title: '設定の更新に失敗しました',
        status: 'error',
        duration: 3000,
        isClosable: true,
      });
    }
  };

  return (
    <Box bg={bgColor} p={8} rounded="lg" shadow="base">
      <VStack spacing={6} as="form" onSubmit={handleSubmit}>
        <FormControl>
          <FormLabel>会社名</FormLabel>
          <Input
            value={formData.companyName}
            onChange={(e) =>
              setFormData({ ...formData, companyName: e.target.value })
            }
          />
        </FormControl>

        <FormControl>
          <FormLabel>メールアドレス</FormLabel>
          <Input
            type="email"
            value={formData.email}
            onChange={(e) =>
              setFormData({ ...formData, email: e.target.value })
            }
          />
        </FormControl>

        <FormControl>
          <FormLabel>言語</FormLabel>
          <Select
            value={formData.language}
            onChange={(e) =>
              setFormData({ ...formData, language: e.target.value })
            }
          >
            <option value="ja">日本語</option>
            <option value="en">English</option>
          </Select>
        </FormControl>

        <FormControl>
          <FormLabel>タイムゾーン</FormLabel>
          <Select
            value={formData.timezone}
            onChange={(e) =>
              setFormData({ ...formData, timezone: e.target.value })
            }
          >
            <option value="Asia/Tokyo">Asia/Tokyo</option>
            <option value="America/New_York">America/New_York</option>
            <option value="Europe/London">Europe/London</option>
          </Select>
        </FormControl>

        <FormControl display="flex" alignItems="center">
          <FormLabel mb="0">通知を有効にする</FormLabel>
          <Switch
            isChecked={formData.notificationEnabled}
            onChange={(e) =>
              setFormData({
                ...formData,
                notificationEnabled: e.target.checked,
              })
            }
          />
        </FormControl>

        <Button type="submit" colorScheme="blue" w="full">
          設定を保存
        </Button>
      </VStack>
    </Box>
  );
} 