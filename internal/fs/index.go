package fs

//func (arc *FileSystem) NewFileIndex() error {
//	//if !sem.TryAcquire(1) {
//	//	return nil
//	//}
//	//defer sem.Release(1)
//	//arc.fiMux.Lock()
//	//defer arc.fiMux.Unlock()
//	bin, err := arc.fileIndex.GenerateBinary()
//	if err != nil {
//		return err
//	}
//	pw := bytes.Buffer{}
//
//	zw, err := gzip.NewWriterLevel(&pw, 2)
//	if err != nil {
//		return err
//	}
//
//	_, err = zw.Write(bin)
//	if err != nil {
//		return err
//	}
//	// bin = pw.Bytes()
//	// bin = snappy.Encode([]byte{}, bin)
//	err = zw.Close()
//	if err != nil {
//		return err
//	}
//	bin = pw.Bytes()
//
//	if arc.fiChecksum == crc32.ChecksumIEEE(bin) {
//		return nil
//	}
//	// for mirroring
//	FileIndexStartBlock := FileIndexStartBlockFlip
//	if arc.fileIndexIsFlip {
//		FileIndexStartBlock = FileIndexStartBlockFlop
//	}
//
//	blockLengthToRead := BLOCKSIZEUSABLE
//	blockNum := int(math.Ceil(float64(len(bin)) / float64(BLOCKSIZEUSABLE)))
//	for i := 0; i < blockNum; i++ {
//		if len(bin[uint64(i)*BLOCKSIZEUSABLE:]) < BLOCKSIZEUSABLE {
//			blockLengthToRead = len(bin[uint64(i)*BLOCKSIZEUSABLE:])
//		}
//		prevBlock := uint32(0)
//		if i > FileIndexStartBlock {
//			prevBlock = uint32(FileIndexStartBlock + i)
//		}
//		var binP []byte
//		if i == blockNum-1 {
//			binP = bin[uint64(i)*BLOCKSIZEUSABLE:]
//		} else {
//			binP = bin[uint64(i)*BLOCKSIZEUSABLE : (uint64(i)*BLOCKSIZEUSABLE)+uint64(blockLengthToRead)]
//		}
//		d := arc.PrepareBlock(binP,
//			uint32(FileIndexStartBlock+i),
//			prevBlock,
//			0,
//			i == blockNum-1)
//		_, err = arc.WriteBlockAt(d, uint32(FileIndexStartBlock+i))
//		if err != nil {
//			arc.log.Error(err)
//			return err
//		}
//	}
//	arc.fiChecksum = crc32.ChecksumIEEE(bin)
//	if arc.fileIndexIsFlip {
//		arc.blockAllocationMap.Remove(FileIndexStartBlockFlip)
//	} else {
//		arc.blockAllocationMap.Remove(FileIndexStartBlockFlop)
//	}
//	arc.fileIndexIsFlip = !arc.fileIndexIsFlip
//
//	return nil
//}
