/**
 * Script para gerar ícones PWA a partir de SVG
 * Requer sharp: npm install -D sharp
 */

import sharp from 'sharp'
import { readFileSync } from 'fs'
import { join } from 'path'

const sizes = [192, 512]
const publicDir = join(process.cwd(), 'public')

async function generateIcons() {
  console.log('Gerando ícones PWA...')
  
  for (const size of sizes) {
    const svgPath = join(publicDir, `pwa-${size}x${size}.svg`)
    const pngPath = join(publicDir, `pwa-${size}x${size}.png`)
    
    try {
      const svg = readFileSync(svgPath)
      await sharp(svg)
        .resize(size, size)
        .png()
        .toFile(pngPath)
      
      console.log(`✅ Gerado: ${pngPath}`)
    } catch (error) {
      console.error(`❌ Erro ao gerar ${pngPath}:`, error.message)
    }
  }
  
  // Gerar apple-touch-icon (180x180)
  try {
    const svg512 = readFileSync(join(publicDir, 'pwa-512x512.svg'))
    await sharp(svg512)
      .resize(180, 180)
      .png()
      .toFile(join(publicDir, 'apple-touch-icon.png'))
    
    console.log('✅ Gerado: apple-touch-icon.png')
  } catch (error) {
    console.error('❌ Erro ao gerar apple-touch-icon.png:', error.message)
  }
  
  // Gerar mask-icon (SVG para Safari)
  try {
    const svg192 = readFileSync(join(publicDir, 'pwa-192x192.svg'))
    await sharp(svg192)
      .resize(512, 512)
      .png()
      .toFile(join(publicDir, 'mask-icon.png'))
    
    console.log('✅ Gerado: mask-icon.png')
  } catch (error) {
    console.error('❌ Erro ao gerar mask-icon.png:', error.message)
  }
  
  console.log('✨ Ícones PWA gerados com sucesso!')
}

generateIcons().catch(console.error)

