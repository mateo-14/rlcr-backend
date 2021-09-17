import { Request, Response } from 'express';
import { client } from '../ds';
import { addOrUpdate, getAll, getByID, getData } from '../services/users';

export function getUser(req: Request, res: Response) {
  client.users
    .fetch(req.userID)
    .then(async (user) => {
      await addOrUpdate({ ...user, ip: req.headers['x-forwarded-for'] || req.socket.remoteAddress } as User);
      const userData = await getData(req.userID);
      res.json({
        avatar: user.avatar,
        username: `${user.username}#${user.discriminator}`,
        id: user.id,
        isAdmin: userData?.isAdmin,
      });
    })
    .catch((err) => {
      console.error(err);
      res.sendStatus(401);
    });
}

export function getAllUsers(req: Request, res: Response) {
  getAll()
    .then((users) => res.json(users))
    .catch((err) => {
      console.error(err);
      res.sendStatus(500);
    });
}

export function getUserByID(req: Request, res: Response) {
  getByID(req.params.id)
    .then((user) => res.json(user))
    .catch((err) => {
      console.error(err);
      res.sendStatus(500);
    });
}
