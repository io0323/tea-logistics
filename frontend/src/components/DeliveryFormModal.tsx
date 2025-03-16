import {
  Modal,
  ModalOverlay,
  ModalContent,
  ModalHeader,
  ModalFooter,
  ModalBody,
  ModalCloseButton,
  Button,
  FormControl,
  FormLabel,
  Input,
  Select,
  Textarea,
  VStack,
  useToast,
} from '@chakra-ui/react';
import { useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import * as yup from 'yup';
import { Delivery, DeliveryStatus, CreateDeliveryRequest, UpdateDeliveryRequest } from '@/types/delivery';

const schema = yup.object().shape({
  orderId: yup.number().required('注文IDは必須です'),
  customerName: yup.string().required('顧客名は必須です'),
  customerAddress: yup.string().required('配送先住所は必須です'),
  customerPhone: yup.string().required('電話番号は必須です'),
  estimatedDeliveryDate: yup.string(),
  note: yup.string(),
});

interface DeliveryFormModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (data: CreateDeliveryRequest | UpdateDeliveryRequest) => Promise<void>;
  initialData?: Delivery;
  title: string;
}

/**
 * 配送フォームのモーダルコンポーネント
 */
export default function DeliveryFormModal({
  isOpen,
  onClose,
  onSubmit,
  initialData,
  title,
}: DeliveryFormModalProps) {
  const toast = useToast();
  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<CreateDeliveryRequest>({
    resolver: yupResolver(schema),
    defaultValues: initialData,
  });

  const handleFormSubmit = async (data: CreateDeliveryRequest) => {
    try {
      await onSubmit(data);
      toast({
        title: '保存しました',
        status: 'success',
        duration: 3000,
        isClosable: true,
      });
      onClose();
    } catch (error) {
      toast({
        title: '保存に失敗しました',
        status: 'error',
        duration: 3000,
        isClosable: true,
      });
    }
  };

  return (
    <Modal isOpen={isOpen} onClose={onClose}>
      <ModalOverlay />
      <ModalContent>
        <ModalHeader>{title}</ModalHeader>
        <ModalCloseButton />
        <ModalBody>
          <form onSubmit={handleSubmit(handleFormSubmit)}>
            <VStack spacing={4}>
              <FormControl isInvalid={!!errors.orderId}>
                <FormLabel>注文ID</FormLabel>
                <Input
                  type="number"
                  {...register('orderId')}
                />
                <FormLabel color="red.500">
                  {errors.orderId?.message}
                </FormLabel>
              </FormControl>

              <FormControl isInvalid={!!errors.customerName}>
                <FormLabel>顧客名</FormLabel>
                <Input {...register('customerName')} />
                <FormLabel color="red.500">
                  {errors.customerName?.message}
                </FormLabel>
              </FormControl>

              <FormControl isInvalid={!!errors.customerAddress}>
                <FormLabel>配送先住所</FormLabel>
                <Input {...register('customerAddress')} />
                <FormLabel color="red.500">
                  {errors.customerAddress?.message}
                </FormLabel>
              </FormControl>

              <FormControl isInvalid={!!errors.customerPhone}>
                <FormLabel>電話番号</FormLabel>
                <Input {...register('customerPhone')} />
                <FormLabel color="red.500">
                  {errors.customerPhone?.message}
                </FormLabel>
              </FormControl>

              <FormControl>
                <FormLabel>予定配送日</FormLabel>
                <Input
                  type="date"
                  {...register('estimatedDeliveryDate')}
                />
              </FormControl>

              <FormControl>
                <FormLabel>備考</FormLabel>
                <Textarea {...register('note')} />
              </FormControl>

              <ModalFooter>
                <Button variant="ghost" mr={3} onClick={onClose}>
                  キャンセル
                </Button>
                <Button
                  colorScheme="blue"
                  type="submit"
                  isLoading={isSubmitting}
                >
                  保存
                </Button>
              </ModalFooter>
            </VStack>
          </form>
        </ModalBody>
      </ModalContent>
    </Modal>
  );
} 