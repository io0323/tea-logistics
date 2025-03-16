import { HttpException, HttpStatus } from '@nestjs/common';

export class UserAlreadyExistsException extends HttpException {
  constructor(email: string) {
    super(
      {
        statusCode: HttpStatus.CONFLICT,
        message: `メールアドレス ${email} は既に登録されています`,
      },
      HttpStatus.CONFLICT,
    );
  }
}

export class InvalidCredentialsException extends HttpException {
  constructor() {
    super(
      {
        statusCode: HttpStatus.UNAUTHORIZED,
        message: 'メールアドレスまたはパスワードが正しくありません',
      },
      HttpStatus.UNAUTHORIZED,
    );
  }
}

export class UserNotFoundException extends HttpException {
  constructor(userId: number) {
    super(
      {
        statusCode: HttpStatus.NOT_FOUND,
        message: `ユーザーID ${userId} が見つかりません`,
      },
      HttpStatus.NOT_FOUND,
    );
  }
} 