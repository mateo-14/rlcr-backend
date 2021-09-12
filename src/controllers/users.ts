import { Request, Response } from 'express';
import { client } from '../ds';
import { addOrUpdate, getAll, getData } from '../services/users';

export const getUser = (req: Request, res: Response) => {
  client.users
    .fetch(req.userID)
    .then(async (user) => {
      await addOrUpdate(user);
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
};

export const getAllUsers = (req: Request, res: Response) => {
  getAll()
    .then((users) => res.json(users))
    .catch((err) => {
      console.error(err);
      res.sendStatus(500);
    });
};
