import express from 'express';
import cors from 'cors';
import cookieParser from 'cookie-parser';
import authRouter from './routes/auth';
import usersRouter from './routes/users';
import settingsRouter from './routes/settings';
import ordersRouter from './routes/orders';

const app = express();

app.use(cors({ credentials: true, origin: 'http://localhost:3000' }));
app.use(express.json());
app.use(cookieParser());

app.use('/auth', authRouter);
app.use('/users', usersRouter);
app.use('/settings', settingsRouter);
app.use('/orders', ordersRouter);

export default app;
