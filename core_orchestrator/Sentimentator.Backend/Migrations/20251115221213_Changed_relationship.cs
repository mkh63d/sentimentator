using System;
using Microsoft.EntityFrameworkCore.Migrations;

#nullable disable

namespace Sentimentator.Backend.Migrations
{
    /// <inheritdoc />
    public partial class Changed_relationship : Migration
    {
        /// <inheritdoc />
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropForeignKey(
                name: "FK_Analyses_AspNetUsers_UserId",
                table: "Analyses");

            migrationBuilder.AlterColumn<Guid>(
                name: "UserId",
                table: "Analyses",
                type: "uuid",
                nullable: true,
                oldClrType: typeof(Guid),
                oldType: "uuid");

            migrationBuilder.AddForeignKey(
                name: "FK_Analyses_AspNetUsers_UserId",
                table: "Analyses",
                column: "UserId",
                principalTable: "AspNetUsers",
                principalColumn: "Id");
        }

        /// <inheritdoc />
        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropForeignKey(
                name: "FK_Analyses_AspNetUsers_UserId",
                table: "Analyses");

            migrationBuilder.AlterColumn<Guid>(
                name: "UserId",
                table: "Analyses",
                type: "uuid",
                nullable: false,
                defaultValue: new Guid("00000000-0000-0000-0000-000000000000"),
                oldClrType: typeof(Guid),
                oldType: "uuid",
                oldNullable: true);

            migrationBuilder.AddForeignKey(
                name: "FK_Analyses_AspNetUsers_UserId",
                table: "Analyses",
                column: "UserId",
                principalTable: "AspNetUsers",
                principalColumn: "Id",
                onDelete: ReferentialAction.Cascade);
        }
    }
}
