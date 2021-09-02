import dotenv from 'dotenv';
dotenv.config();
import app from './app';
import { client } from './ds';

const PORT = process.env.PORT || 8080;
app.listen(PORT, () => {
  console.log(`App listen on port ${PORT}`);
});

client.login(process.env.CLIENT_TOKEN).catch((err) => console.error(err));
