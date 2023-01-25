package util

import (
	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
	"os"
	"runtime"
	"syscall"
	"unsafe"
)

const runType = runtime.GOOS

// mmap 是可以实现用户内存空间（不占用物理内存的虚拟内存）和内核空间的映射（是一种零拷贝技术）
// 用户可以使用指针操作内存，而不需要进行write()和read()的系统调用
// 系统会定时将内存中的数据同步到磁盘文件中

// 用户进程调用 mmap()
// ，从用户态陷入内核态，将内核缓冲区映射到用户缓存区；
// DMA 控制器将数据从硬盘拷贝到内核缓冲区（可见其使用了 Page Cache 机制）；
// mmap()
// 返回，上下文从内核态切换回用户态；
// 用户进程调用 write()
// ，尝试把文件数据写到内核里的套接字缓冲区（在这里处于一种懒加载的模式，并不会马上分配物理内存，当用户使用类似write()的系统调用是，cpu会陷入缺页异常从而进行文件的加载），再次陷入内核态；
// CPU 将内核缓冲区中的数据拷贝到的套接字缓冲区；
// DMA 控制器将数据从套接字缓冲区拷贝到网卡完成数据传输；
// write()
// 返回，上下文从内核态切换回用户态。

// mmap的懒加载导致的缺页异常会导致性能的下降，因为cpu需要加载文件
func standMmap(fd *os.File, size int64) ([]byte, error) {
	return syscall.Mmap(int(fd.Fd()), 0, int(size), syscall.PROT_WRITE, syscall.MAP_SHARED)
}

// MADV_NORMAL ：这是默认行为。分页是以簇的形式（较小的一个系统分页大小的整数倍）传输的。这个值会导致一些预先读和事后读。
// MADV_RANDOM ：这个区域中的分页会被随机访问，这样预先读将不会带来任何好处，因此内核在每次读取时所取出的数据量应该尽可能少。
// MADV_SEQUENTIAL ：在这个范围中的分页只会被访问一次，并且是顺序访问，因此内核可以激进地预先读，并且分页在被访问之后就可以将其释放了
// MADV_WILLNEED ：预先读取这个区域中的分页以备将来的访问之需。MADV_WILLNEED 操作的效果与 Linux特有的 readahead()系统调用和 posix_fadvise() POSIX_FADV_WILLNEED 操作的效果类似
// MADV_DONTNEED ：调用进程不再要求这个区域中的分页驻留在内存中。这个标记的精确效果在不同 UNIX 实现上是不同的。

// darwin 不支持内存建议操作，因此由外部保证使用选择
func standMadvise(byteData []byte) error {
	advice := unix.MADV_NORMAL
	_, _, eno := syscall.Syscall(syscall.SYS_MADVISE, uintptr(unsafe.Pointer(&byteData[0])),
		uintptr(len(byteData)), uintptr(advice))
	if eno != 0 {
		return eno
	}
	return nil
}

type MmapFile struct {
	Fd       *os.File
	ByteData []byte
}

func mmap(fd *os.File, sz int) (*MmapFile, error) {
	logger := Logger().SetField("mmap")
	filename := fd.Name()
	fi, err := fd.Stat()
	if err != nil {
		return nil, errors.Wrapf(err, "cannot stat file: %s err: %v", filename, err)
	}

	fileSize := fi.Size()
	if sz > 0 && fileSize == 0 {
		// If file is empty, truncate it to sz.
		if err := fd.Truncate(int64(sz)); err != nil {
			return nil, errors.Wrapf(err, "error while truncation")
		}
		fileSize = int64(sz)
	}

	bytes, err := standMmap(fd, fileSize)
	if err != nil {
		return nil, errors.Wrapf(err, "mapping failed: %s err: %v", filename, err)
	}
	if runType == "linux" { // 对于linux系统，可以进行内存建议操作
		err = standMadvise(bytes)
		if err != nil {
			logger.Warn("madvise failed: %s err: %v", filename, err)
		}
	}

	return &MmapFile{Fd: fd, ByteData: bytes}, nil
}

func NewMmapFile(filename string) (*MmapFile, error) {
	fd, err := os.OpenFile(filename, os.O_RDWR, 0666)
	if err != nil {
		return nil, errors.Wrapf(err, "open file failed: %s err: %v", filename, err)
	}
	return mmap(fd, 0)
}
