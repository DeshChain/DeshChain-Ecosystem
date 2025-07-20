/*
Copyright 2024 DeshChain Foundation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import React, { useEffect, useRef, useState } from 'react';
import {
  Box,
  Paper,
  Typography,
  Button,
  IconButton,
  Tooltip,
  CircularProgress,
  useTheme
} from '@mui/material';
import {
  Download as DownloadIcon,
  Share as ShareIcon,
  Print as PrintIcon,
  ContentCopy as CopyIcon,
  QrCode as QrCodeIcon
} from '@mui/icons-material';
import QRCode from 'qrcode';
import { motion } from 'framer-motion';
import toast from 'react-hot-toast';

interface QRCodeData {
  receiptId: string;
  orderId: string;
  amount: string;
  currency: string;
  sender: string;
  receiver: string;
  timestamp: string;
  verificationCode: string;
  culturalQuote?: string;
  festivalBonus?: string;
}

interface QRCodeGeneratorProps {
  data: QRCodeData;
  size?: number;
  showActions?: boolean;
  format?: 'canvas' | 'svg' | 'image';
  errorCorrectionLevel?: 'L' | 'M' | 'Q' | 'H';
  includeText?: boolean;
  culturalDesign?: boolean;
}

export const QRCodeGenerator: React.FC<QRCodeGeneratorProps> = ({
  data,
  size = 256,
  showActions = true,
  format = 'canvas',
  errorCorrectionLevel = 'H',
  includeText = true,
  culturalDesign = true
}) => {
  const theme = useTheme();
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const [isGenerating, setIsGenerating] = useState(true);
  const [qrDataUrl, setQrDataUrl] = useState<string>('');

  useEffect(() => {
    generateQRCode();
  }, [data, size, errorCorrectionLevel, culturalDesign]);

  const generateQRCode = async () => {
    setIsGenerating(true);

    try {
      // Create QR data payload
      const qrPayload = {
        v: '1.0', // Version
        rid: data.receiptId,
        oid: data.orderId,
        amt: data.amount,
        cur: data.currency,
        sndr: data.sender.slice(0, 20), // Truncate for QR size
        rcvr: data.receiver.slice(0, 20),
        ts: data.timestamp,
        vc: data.verificationCode,
        ...(data.culturalQuote && { cq: data.culturalQuote.slice(0, 50) }),
        ...(data.festivalBonus && { fb: data.festivalBonus })
      };

      const qrString = JSON.stringify(qrPayload);

      if (format === 'canvas' && canvasRef.current) {
        const canvas = canvasRef.current;
        
        // Generate QR code
        await QRCode.toCanvas(canvas, qrString, {
          width: size,
          errorCorrectionLevel,
          margin: culturalDesign ? 4 : 2,
          color: {
            dark: culturalDesign ? theme.palette.primary.main : '#000000',
            light: '#FFFFFF'
          }
        });

        // Add cultural design elements if enabled
        if (culturalDesign) {
          const ctx = canvas.getContext('2d');
          if (ctx) {
            // Add border pattern
            ctx.strokeStyle = theme.palette.secondary.main;
            ctx.lineWidth = 3;
            ctx.setLineDash([5, 5]);
            ctx.strokeRect(10, 10, size - 20, size - 20);
            
            // Add corner decorations
            drawCornerDesigns(ctx, size);
          }
        }

        // Generate data URL for download/share
        setQrDataUrl(canvas.toDataURL('image/png'));
      } else {
        // Generate as data URL directly
        const dataUrl = await QRCode.toDataURL(qrString, {
          width: size,
          errorCorrectionLevel,
          margin: culturalDesign ? 4 : 2,
          color: {
            dark: culturalDesign ? theme.palette.primary.main : '#000000',
            light: '#FFFFFF'
          },
          type: 'image/png'
        });
        setQrDataUrl(dataUrl);
      }

      setIsGenerating(false);
    } catch (error) {
      console.error('Failed to generate QR code:', error);
      setIsGenerating(false);
      toast.error('Failed to generate QR code');
    }
  };

  const drawCornerDesigns = (ctx: CanvasRenderingContext2D, size: number) => {
    const cornerSize = 20;
    
    // Top-left corner
    ctx.beginPath();
    ctx.moveTo(10, 10 + cornerSize);
    ctx.lineTo(10, 10);
    ctx.lineTo(10 + cornerSize, 10);
    ctx.strokeStyle = theme.palette.primary.main;
    ctx.lineWidth = 3;
    ctx.setLineDash([]);
    ctx.stroke();

    // Top-right corner
    ctx.beginPath();
    ctx.moveTo(size - 10 - cornerSize, 10);
    ctx.lineTo(size - 10, 10);
    ctx.lineTo(size - 10, 10 + cornerSize);
    ctx.stroke();

    // Bottom-left corner
    ctx.beginPath();
    ctx.moveTo(10, size - 10 - cornerSize);
    ctx.lineTo(10, size - 10);
    ctx.lineTo(10 + cornerSize, size - 10);
    ctx.stroke();

    // Bottom-right corner
    ctx.beginPath();
    ctx.moveTo(size - 10 - cornerSize, size - 10);
    ctx.lineTo(size - 10, size - 10);
    ctx.lineTo(size - 10, size - 10 - cornerSize);
    ctx.stroke();
  };

  const handleDownload = () => {
    try {
      const link = document.createElement('a');
      link.download = `receipt-${data.receiptId}.png`;
      link.href = qrDataUrl;
      link.click();
      toast.success('QR code downloaded successfully');
    } catch (error) {
      toast.error('Failed to download QR code');
    }
  };

  const handleCopy = async () => {
    try {
      if (canvasRef.current) {
        canvasRef.current.toBlob(async (blob) => {
          if (blob) {
            await navigator.clipboard.write([
              new ClipboardItem({ 'image/png': blob })
            ]);
            toast.success('QR code copied to clipboard');
          }
        });
      }
    } catch (error) {
      // Fallback to copying the verification URL
      const verificationUrl = `https://deshchain.org/verify/${data.receiptId}`;
      await navigator.clipboard.writeText(verificationUrl);
      toast.success('Verification link copied to clipboard');
    }
  };

  const handleShare = async () => {
    const verificationUrl = `https://deshchain.org/verify/${data.receiptId}`;
    
    if (navigator.share) {
      try {
        // Convert canvas to blob for sharing
        if (canvasRef.current) {
          canvasRef.current.toBlob(async (blob) => {
            if (blob) {
              const file = new File([blob], `receipt-${data.receiptId}.png`, { type: 'image/png' });
              await navigator.share({
                title: 'DeshChain Money Order Receipt',
                text: `Receipt #${data.receiptId} - Amount: ${data.amount} ${data.currency}`,
                files: [file],
                url: verificationUrl
              });
            }
          });
        }
      } catch (error) {
        // Fallback to sharing just the URL
        await navigator.share({
          title: 'DeshChain Money Order Receipt',
          text: `Receipt #${data.receiptId} - Amount: ${data.amount} ${data.currency}`,
          url: verificationUrl
        });
      }
    } else {
      // Fallback to copying the URL
      await navigator.clipboard.writeText(verificationUrl);
      toast.success('Verification link copied to clipboard');
    }
  };

  const handlePrint = () => {
    const printWindow = window.open('', '_blank');
    if (printWindow) {
      printWindow.document.write(`
        <html>
          <head>
            <title>Money Order Receipt - ${data.receiptId}</title>
            <style>
              body {
                font-family: Arial, sans-serif;
                display: flex;
                flex-direction: column;
                align-items: center;
                padding: 20px;
              }
              .receipt-header {
                text-align: center;
                margin-bottom: 20px;
              }
              .qr-container {
                margin: 20px 0;
              }
              .receipt-details {
                margin-top: 20px;
                padding: 20px;
                border: 1px solid #ccc;
                border-radius: 8px;
              }
              .detail-row {
                display: flex;
                justify-content: space-between;
                margin: 10px 0;
              }
              @media print {
                body { margin: 0; }
              }
            </style>
          </head>
          <body>
            <div class="receipt-header">
              <h1>DeshChain Money Order</h1>
              <h2>Receipt #${data.receiptId}</h2>
            </div>
            <div class="qr-container">
              <img src="${qrDataUrl}" alt="QR Code" />
            </div>
            <div class="receipt-details">
              <div class="detail-row">
                <strong>Amount:</strong>
                <span>${data.amount} ${data.currency}</span>
              </div>
              <div class="detail-row">
                <strong>From:</strong>
                <span>${data.sender}</span>
              </div>
              <div class="detail-row">
                <strong>To:</strong>
                <span>${data.receiver}</span>
              </div>
              <div class="detail-row">
                <strong>Date:</strong>
                <span>${new Date(data.timestamp).toLocaleString()}</span>
              </div>
              <div class="detail-row">
                <strong>Verification:</strong>
                <span>${data.verificationCode}</span>
              </div>
            </div>
            ${data.culturalQuote ? `
              <div style="margin-top: 20px; text-align: center; font-style: italic;">
                "${data.culturalQuote}"
              </div>
            ` : ''}
          </body>
        </html>
      `);
      printWindow.document.close();
      printWindow.print();
    }
  };

  return (
    <motion.div
      initial={{ opacity: 0, scale: 0.9 }}
      animate={{ opacity: 1, scale: 1 }}
      transition={{ duration: 0.3 }}
    >
      <Paper
        sx={{
          p: culturalDesign ? 3 : 2,
          textAlign: 'center',
          background: culturalDesign 
            ? 'linear-gradient(135deg, rgba(255,107,53,0.05), rgba(19,136,8,0.05))'
            : 'transparent'
        }}
      >
        {isGenerating ? (
          <Box display="flex" justifyContent="center" alignItems="center" height={size}>
            <CircularProgress />
          </Box>
        ) : (
          <>
            {/* QR Code Display */}
            <Box display="flex" justifyContent="center" mb={2}>
              {format === 'canvas' ? (
                <canvas
                  ref={canvasRef}
                  width={size}
                  height={size}
                  style={{
                    border: culturalDesign ? `2px solid ${theme.palette.divider}` : 'none',
                    borderRadius: culturalDesign ? 8 : 0
                  }}
                />
              ) : (
                <img
                  src={qrDataUrl}
                  alt="QR Code"
                  width={size}
                  height={size}
                  style={{
                    border: culturalDesign ? `2px solid ${theme.palette.divider}` : 'none',
                    borderRadius: culturalDesign ? 8 : 0
                  }}
                />
              )}
            </Box>

            {/* Receipt Info */}
            {includeText && (
              <Box mb={2}>
                <Typography variant="h6" gutterBottom>
                  Receipt #{data.receiptId}
                </Typography>
                <Typography variant="body2" color="text.secondary">
                  {data.amount} {data.currency}
                </Typography>
                <Typography variant="caption" color="text.secondary">
                  {new Date(data.timestamp).toLocaleString()}
                </Typography>
              </Box>
            )}

            {/* Action Buttons */}
            {showActions && (
              <Box display="flex" gap={1} justifyContent="center">
                <Tooltip title="Download QR Code">
                  <IconButton onClick={handleDownload} color="primary">
                    <DownloadIcon />
                  </IconButton>
                </Tooltip>
                <Tooltip title="Copy QR Code">
                  <IconButton onClick={handleCopy} color="primary">
                    <CopyIcon />
                  </IconButton>
                </Tooltip>
                <Tooltip title="Share Receipt">
                  <IconButton onClick={handleShare} color="primary">
                    <ShareIcon />
                  </IconButton>
                </Tooltip>
                <Tooltip title="Print Receipt">
                  <IconButton onClick={handlePrint} color="primary">
                    <PrintIcon />
                  </IconButton>
                </Tooltip>
              </Box>
            )}

            {/* Verification URL */}
            <Box mt={2} p={1} bgcolor="background.default" borderRadius={1}>
              <Typography variant="caption" color="text.secondary">
                Verify at: deshchain.org/verify/{data.receiptId}
              </Typography>
            </Box>
          </>
        )}
      </Paper>
    </motion.div>
  );
};