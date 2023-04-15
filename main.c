#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <libavformat/avformat.h>
#include <libavcodec/avcodec.h>
#include <libavutil/avutil.h>

const size_t MAX_LINES = 10;
const size_t LINE_BUFF = 256;

int process_each(char** lines, size_t n);
int process_one(char* infile);

int main() {
    FILE* fp;
    fp = fopen("test.txt", "r");
    if (fp == NULL) {
        perror("Error opening file");
        return 1;
    }
    
    char buf[LINE_BUFF];
    char** lines = calloc(MAX_LINES, sizeof(char*));

    size_t i = 0;
    while (fgets(buf, LINE_BUFF, fp) && i < MAX_LINES) {
        lines[i] = strdup(buf);
        memset(buf, '\0', LINE_BUFF);
        i++;
    }

    if (ferror(fp) || i == 0) {
        perror("Error reading file");
        fclose(fp);
        return 1;
    }
    fclose(fp);

    process_one(lines[0]);

    for (size_t j = 0; j < i; j++) {
        free(lines[j]);
    }
    free(lines);

    return 0;
}

int process_each(char** lines, size_t n) {
    for (int i = 0; i < n; i++) {
        printf("Processing file: %s\n", lines[i]);
    }
    return 0;
}

int process_one(char* infile) {
    AVFormatContext* input_ctx = NULL;
    AVFormatContext* output_ctx = NULL;
    AVPacket packet;
    const char* outfile = "output.mp4";

    infile[strcspn(infile, "\n")] = 0;
    int ret = avformat_open_input(&input_ctx, infile, NULL, NULL);
    if (ret < 0) {
        fprintf(stderr, "Failed to open input file: %s\n", infile);
        return 1;
    }

    ret = avformat_find_stream_info(input_ctx, NULL);
    if (ret < 0) {
        fprintf(stderr, "Failed to find stream info\n");
        return 1;
    }

    ret = avformat_alloc_output_context2(&output_ctx, NULL, NULL, outfile);
    if (ret < 0) {
        fprintf(stderr, "Failed to allocate output context\n");
        return 1;
    }

    ret = avio_open(&output_ctx->pb, outfile, AVIO_FLAG_WRITE);
    if (ret < 0) {
        fprintf(stderr, "Failed to open output file\n");
        return 1;
    }

    for (int i = 0; i < input_ctx->nb_streams; i++) {
        AVStream* in_stream = input_ctx->streams[i];
        AVStream* out_stream = avformat_new_stream(output_ctx, NULL);
        if (!out_stream) {
            fprintf(stderr, "Failed allocating output stream\n");
            return 1;
        }
        ret = avcodec_parameters_copy(out_stream->codecpar, in_stream->codecpar);
        if (ret < 0) {
            fprintf(stderr, "Failed to copy codec parameters\n");
            return 1;
        }
    }

    ret = avformat_write_header(output_ctx, NULL);
    if (ret < 0) {
        fprintf(stderr, "Failed to write header\n");
        return 1;
    }

    while (av_read_frame(input_ctx, &packet) >= 0) {
        AVStream* in_stream = input_ctx->streams[packet.stream_index];
        AVStream* out_stream = output_ctx->streams[packet.stream_index];
        packet.pts = av_rescale_q_rnd(packet.pts, in_stream->time_base, out_stream->time_base, AV_ROUND_NEAR_INF);
        packet.dts = av_rescale_q_rnd(packet.dts, in_stream->time_base, out_stream->time_base, AV_ROUND_NEAR_INF);
        packet.duration = av_rescale_q(packet.duration, in_stream->time_base, out_stream->time_base);
        packet.pos = -1;
        ret = av_interleaved_write_frame(output_ctx, &packet);
        if (ret < 0) {
            fprintf(stderr, "Failed to write frame (packet)\n");
            break;
        }
        av_packet_unref(&packet);
    }

    av_write_trailer(output_ctx);

    avformat_close_input(&input_ctx);
    avformat_free_context(input_ctx);
    avformat_free_context(output_ctx);

    return 0;
}
