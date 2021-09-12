import { PaymentMethodID, OrderStatus, OrderMode } from './index';

declare global {
  declare namespace Express {
    interface Request {
      userID: string;
    }
  }

  export interface PaymentMethod {
    id: PaymentMethodID;
    name: string;
  }

  export interface Settings {
    buyEnabled: boolean;
    sellEnabled: boolean;
    creditBuyValue: number;
    creditSellValue: number;
    maxBuy: number;
    maxSell: number;
    paymentMethods: PaymentMethod[];
  }

  export interface Order {
    id: string;
    userID: string;
    credits: number;
    price: number;
    paymentMethodID: PaymentMethodID;
    mode: OrderMode;
    dni?: string;
    account: string;
    paymentAccount?: string;
    cvuAlias?: string;
    createdAt: number;
    status: OrderStatus;
  }
}
