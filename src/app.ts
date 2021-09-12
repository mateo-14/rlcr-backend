import express from 'express';
import cors from 'cors';
import cookieParser from 'cookie-parser';
import authRouter from './routes/auth';
import usersRouter from './routes/users';
import settingsRouter from './routes/settings';
import ordersRouter from './routes/orders';

const app = express();

const whitelist = process.env.CORS_WHITELIST?.split(', ');
app.use(
  cors({
    credentials: true,
    origin: (origin, callback) =>
      !origin
        ? callback(null, true)
        : whitelist?.indexOf(origin) !== -1
        ? callback(null, true)
        : callback(new Error('Not allowed by CORS')),
  })
);

app.use(express.json());
app.use(cookieParser());

app.use('/auth', authRouter);
app.use('/users', usersRouter);
app.use('/settings', settingsRouter);
app.use('/orders', ordersRouter);

export default app;
