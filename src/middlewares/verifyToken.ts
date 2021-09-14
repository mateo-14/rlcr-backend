import { NextFunction, Request, Response } from 'express';
import { verify } from '../services/tokens';

export default function verifyToken(req: Request, res: Response, next: NextFunction) {
  if (req.cookies?.token) {
    return verify(req.cookies.token)
      .then((userID) => {
        req.userID = userID;
        next();
      })
      .catch((err) => {
        res.cookie('token', '', { httpOnly: true, expires: new Date(0) });
        console.error(err);
        res.sendStatus(401);
      });
  }
  res.sendStatus(401);
}
