@dnl = internal constant [4 x i8] c"%d\0A\00"
declare i32 @printf(i8*, ...)
declare i32 @puts(i8*)
define void @printInt(i32 %x) {
   %t0 = getelementptr [4 x i8], [4 x i8]* @dnl, i32 0, i32 0
   call i32 (i8*, ...) @printf(i8* %t0, i32 %x)
   ret void
}
define i32 @main() {
   %1 = sdiv i32 8, 2
   %2 = add i32 %1, 2
   call void @printInt(i32 %2)
	ret i32 0
}