import { OrderMode } from '../../@types';
import firestore from '../firestore';

let cachedSettings: Settings | undefined;

const getSettings = (): Promise<Settings> => {
  return new Promise(async (resolve, reject) => {
    if (cachedSettings) return resolve(cachedSettings);
    try {
      const docSnap = await firestore.collection('settings').doc('default').get();
      if (docSnap.exists && docSnap.data()) {
        cachedSettings = docSnap.data() as Settings;
        return resolve(cachedSettings);
      }
    } catch (err) {
      return reject(err);
    }
    reject('Settings does not exists');
  });
};

export const sanitizeCredits = async (credits: number, mode: OrderMode) => {
  const settings = await getSettings();
  const max = mode === OrderMode.Buy ? settings?.maxBuy : settings?.maxSell;
  return Math.max(100, Math.min(Math.round(credits / 10) * 10, max));
};

export default { getSettings, sanitizeCredits };
